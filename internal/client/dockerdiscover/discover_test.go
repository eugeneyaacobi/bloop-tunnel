package dockerdiscover

import (
	"context"
	"io"
	"net/http"
	"strings"
	"testing"
)

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(r *http.Request) (*http.Response, error) { return f(r) }

func TestDiscoverFiltersInfrastructureAndSortsCandidates(t *testing.T) {
	client := &Client{
		SocketPath: "/tmp/not-used.sock",
		HTTPClient: &http.Client{Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
			body := `[
				{"Id":"1","Names":["/frontend"],"Image":"frontend:latest","Ports":[{"PrivatePort":3000,"Type":"tcp"}]},
				{"Id":"2","Names":["/postgres"],"Image":"postgres:16","Ports":[{"PrivatePort":5432,"Type":"tcp"}]},
				{"Id":"3","Names":["/api"],"Image":"api:latest","Ports":[{"PrivatePort":8080,"Type":"tcp"}]}
			]`
			return &http.Response{StatusCode: http.StatusOK, Body: io.NopCloser(strings.NewReader(body)), Header: make(http.Header)}, nil
		})},
	}

	// bypass os.Stat through an existing path.
	client.SocketPath = "/dev/null"
	candidates, err := client.Discover(context.Background())
	if err != nil {
		t.Fatalf("Discover returned error: %v", err)
	}
	if len(candidates) != 2 {
		t.Fatalf("candidate count = %d (%#v)", len(candidates), candidates)
	}
	if candidates[0].ContainerName != "api" || candidates[1].ContainerName != "frontend" {
		t.Fatalf("unexpected candidate order: %#v", candidates)
	}
}
