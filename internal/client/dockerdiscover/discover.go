package dockerdiscover

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"time"
)

const DefaultSocketPath = "/var/run/docker.sock"

type Candidate struct {
	ContainerID   string
	ContainerName string
	Image         string
	Port          int
	LocalAddr     string
	SuggestedName string
}

type containerSummary struct {
	ID    string   `json:"Id"`
	Names []string `json:"Names"`
	Image string   `json:"Image"`
	Ports []struct {
		PrivatePort int    `json:"PrivatePort"`
		PublicPort  int    `json:"PublicPort"`
		Type        string `json:"Type"`
	} `json:"Ports"`
}

type Discoverer interface {
	Discover(context.Context) ([]Candidate, error)
}

type Client struct {
	SocketPath string
	HTTPClient *http.Client
}

func NewClient(socketPath string) *Client {
	if socketPath == "" {
		socketPath = DefaultSocketPath
	}
	transport := &http.Transport{
		DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
			var d net.Dialer
			return d.DialContext(ctx, "unix", socketPath)
		},
	}
	return &Client{
		SocketPath: socketPath,
		HTTPClient: &http.Client{Transport: transport, Timeout: 4 * time.Second},
	}
}

func IsSupported() bool {
	return runtime.GOOS == "linux" || runtime.GOOS == "darwin"
}

func (c *Client) Discover(ctx context.Context) ([]Candidate, error) {
	if !IsSupported() {
		return nil, errors.New("docker discovery is unsupported on this platform")
	}
	if c.SocketPath == "" {
		c.SocketPath = DefaultSocketPath
	}
	if _, err := os.Stat(c.SocketPath); err != nil {
		return nil, fmt.Errorf("docker socket unavailable: %w", err)
	}
	if c.HTTPClient == nil {
		c.HTTPClient = NewClient(c.SocketPath).HTTPClient
	}

	reqURL := (&url.URL{Scheme: "http", Host: "docker", Path: "/containers/json"}).String()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, reqURL, nil)
	if err != nil {
		return nil, err
	}
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("docker daemon unavailable: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("docker discovery failed: status %d", resp.StatusCode)
	}

	var containers []containerSummary
	if err := json.NewDecoder(resp.Body).Decode(&containers); err != nil {
		return nil, fmt.Errorf("decode docker response: %w", err)
	}

	candidates := make([]Candidate, 0)
	for _, container := range containers {
		name := canonicalContainerName(container.Names)
		if excluded(name, container.Image) {
			continue
		}
		for _, port := range container.Ports {
			if port.Type != "tcp" || port.PrivatePort == 0 || !looksHTTPish(port.PrivatePort) {
				continue
			}
			candidates = append(candidates, Candidate{
				ContainerID:   container.ID,
				ContainerName: name,
				Image:         container.Image,
				Port:          port.PrivatePort,
				LocalAddr:     net.JoinHostPort(name, strconv.Itoa(port.PrivatePort)),
				SuggestedName: suggestedName(name, port.PrivatePort),
			})
		}
	}

	sort.Slice(candidates, func(i, j int) bool {
		if candidates[i].ContainerName == candidates[j].ContainerName {
			return candidates[i].Port < candidates[j].Port
		}
		return candidates[i].ContainerName < candidates[j].ContainerName
	})
	return candidates, nil
}

func canonicalContainerName(names []string) string {
	for _, name := range names {
		name = strings.TrimPrefix(name, "/")
		if name != "" {
			return filepath.Base(name)
		}
	}
	return "container"
}

func suggestedName(name string, port int) string {
	clean := strings.Map(func(r rune) rune {
		switch {
		case r >= 'a' && r <= 'z', r >= 'A' && r <= 'Z', r >= '0' && r <= '9':
			return r
		case r == '-', r == '_':
			return '-'
		default:
			return '-'
		}
	}, strings.TrimSpace(name))
	clean = strings.Trim(clean, "-")
	if clean == "" {
		clean = "app"
	}
	return fmt.Sprintf("%s-%d", strings.ToLower(clean), port)
}

