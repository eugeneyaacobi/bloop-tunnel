package setup

import (
	"bufio"
	"bytes"
	"context"
	"strings"
	"testing"

	"bloop-tunnel/internal/client/dockerdiscover"
	"bloop-tunnel/internal/config"
)

func TestDockerDiscoverySkippedWhenNotRequested(t *testing.T) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	input := strings.NewReader(strings.Join([]string{
		"", "", "", "", "", "", // defaults (no discovery prompt)
		"done",
	}, "\n") + "\n")

	err := Run(&stdout, &stderr, Options{
		ConfigPath:     "-",
		OutputMode:     OutputYAML,
		Stdin:          input,
		DiscoverDocker: false,
		Discoverer: fakeDiscoverer{candidates: []dockerdiscover.Candidate{
			{ContainerName: "app", Image: "app:latest", LocalAddr: "app:3000", SuggestedName: "app-3000"},
		}},
	})
	if err != nil {
		t.Fatalf("Run returned error: %v", err)
	}
	if strings.Contains(stdout.String(), "app-3000") {
		t.Fatal("discovered tunnel should not appear when discovery is not requested")
	}
}

func TestDockerDiscoveryNoCandidatesContinuesToManualEntry(t *testing.T) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	input := strings.NewReader(strings.Join([]string{
		"", "", "", "", "", "", // defaults
		"y",    // attempt discovery
		"done", // finish tunnel editing (default tunnel remains)
	}, "\n") + "\n")

	err := Run(&stdout, &stderr, Options{
		ConfigPath:     "-",
		OutputMode:     OutputYAML,
		Stdin:          input,
		DiscoverDocker: true,
		Discoverer:     fakeDiscoverer{candidates: nil},
	})
	if err != nil {
		t.Fatalf("Run returned error: %v", err)
	}
	if !strings.Contains(stdout.String(), "name: app") {
		t.Fatal("default tunnel should remain when no Docker candidates found")
	}
}

func TestDockerDiscoveryErrorHandledGracefully(t *testing.T) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	input := strings.NewReader(strings.Join([]string{
		"", "", "", "", "", "", // defaults
		"y",    // attempt discovery
		"done", // finish
	}, "\n") + "\n")

	err := Run(&stdout, &stderr, Options{
		ConfigPath:     "-",
		OutputMode:     OutputYAML,
		Stdin:          input,
		DiscoverDocker: true,
		Discoverer:     fakeDiscoverer{err: context.DeadlineExceeded},
	})
	if err != nil {
		t.Fatalf("Run returned error: %v", err)
	}
	if !strings.Contains(stderr.String(), "docker discovery unavailable") {
		t.Fatalf("expected graceful discovery warning, got stderr: %q", stderr.String())
	}
}

func TestApplyDockerDiscoveryConfirmationPerCandidate(t *testing.T) {
	var buf bytes.Buffer
	cfg := &config.ClientConfig{}
	opts := Options{
		Discoverer: fakeDiscoverer{candidates: []dockerdiscover.Candidate{
			{ContainerName: "web", Image: "web:latest", LocalAddr: "web:3000", SuggestedName: "web-3000"},
			{ContainerName: "api", Image: "api:latest", LocalAddr: "api:8080", SuggestedName: "api-8080"},
			{ContainerName: "admin", Image: "admin:latest", LocalAddr: "admin:9090", SuggestedName: "admin-9090"},
		}},
	}

	// Confirm first, deny second, confirm third
	input := "y\nn\ny\n"
	reader := bufio.NewReader(strings.NewReader(input))

	err := applyDockerDiscovery(context.Background(), reader, &buf, cfg, opts)
	if err != nil {
		t.Fatalf("applyDockerDiscovery error: %v", err)
	}

	if len(cfg.Tunnels) != 2 {
		t.Fatalf("expected 2 tunnels, got %d: %+v", len(cfg.Tunnels), cfg.Tunnels)
	}
	if cfg.Tunnels[0].Name != "web-3000" {
		t.Errorf("first tunnel name = %q, want web-3000", cfg.Tunnels[0].Name)
	}
	if cfg.Tunnels[1].Name != "admin-9090" {
		t.Errorf("second tunnel name = %q, want admin-9090", cfg.Tunnels[1].Name)
	}

	output := buf.String()
	if !strings.Contains(output, "Docker discovery candidates:") {
		t.Error("expected candidates header in output")
	}
	if !strings.Contains(output, "web:3000") || !strings.Contains(output, "api:8080") || !strings.Contains(output, "admin:9090") {
		t.Error("expected all candidates listed in output")
	}
}

func TestApplyDockerDiscoveryEmptyCandidates(t *testing.T) {
	var buf bytes.Buffer
	cfg := &config.ClientConfig{
		Tunnels: []config.TunnelConfig{{Name: "existing", LocalAddr: "localhost:3000"}},
	}
	opts := Options{
		Discoverer: fakeDiscoverer{candidates: nil},
	}

	err := applyDockerDiscovery(context.Background(), bufio.NewReader(strings.NewReader("")), &buf, cfg, opts)
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	if !strings.Contains(buf.String(), "No eligible Docker services found") {
		t.Error("expected no-candidates message")
	}
	if len(cfg.Tunnels) != 1 {
		t.Fatal("existing tunnels should not be modified when no candidates found")
	}
}
