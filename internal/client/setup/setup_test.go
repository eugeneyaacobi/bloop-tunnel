package setup

import (
	"bytes"
	"context"
	"strings"
	"testing"

	"bloop-tunnel/internal/client/dockerdiscover"
)

type fakeDiscoverer struct {
	candidates []dockerdiscover.Candidate
	err        error
}

func (f fakeDiscoverer) Discover(context.Context) ([]dockerdiscover.Candidate, error) {
	return f.candidates, f.err
}

func TestRunNonInteractiveYAMLScaffoldToStdout(t *testing.T) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	err := Run(&stdout, &stderr, Options{
		ConfigPath:         "-",
		OutputMode:         OutputYAML,
		NonInteractive:     true,
		ControlPlaneURL:    "https://api.bloop.to",
		RelayURL:           "wss://relay.bloop.to/connect",
		AuthTokenEnv:       "BLOOP_CLIENT_TOKEN",
		EnrollmentTokenEnv: "BLOOP_ENROLLMENT_TOKEN",
	})
	if err != nil {
		t.Fatalf("Run returned error: %v", err)
	}
	got := stdout.String()
	for _, needle := range []string{
		"control_plane_url: https://api.bloop.to",
		"relay_url: wss://relay.bloop.to/connect",
		"auth_token_env: BLOOP_CLIENT_TOKEN",
		"enrollment_token_env: BLOOP_ENROLLMENT_TOKEN",
		"name: app",
	} {
		if !strings.Contains(got, needle) {
			t.Fatalf("output missing %q:\n%s", needle, got)
		}
	}
}

func TestRunNonInteractiveEnvFile(t *testing.T) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	err := Run(&stdout, &stderr, Options{OutputMode: OutputEnvFile, NonInteractive: true})
	if err != nil {
		t.Fatalf("Run returned error: %v", err)
	}
	got := stdout.String()
	for _, needle := range []string{
		"BLOOP_CONTROL_PLANE_URL=https://api.bloop.to",
		"BLOOP_RELAY_URL=wss://relay.bloop.to/connect",
		"BLOOP_AUTH_TOKEN_ENV=BLOOP_CLIENT_TOKEN",
		"BLOOP_TUNNELS_0_NAME=app",
		"BLOOP_TUNNELS_0_LOCAL_ADDR=127.0.0.1:3000",
	} {
		if !strings.Contains(got, needle) {
			t.Fatalf("output missing %q:\n%s", needle, got)
		}
	}
}

func TestRunInteractiveEditsExistingConfigAndAddsTunnel(t *testing.T) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	input := strings.NewReader(strings.Join([]string{
		"",                  // keep control plane
		"",                  // keep relay
		"",                  // keep auth env
		"",                  // keep enrollment env
		"",                  // keep reconnect initial
		"",                  // keep reconnect max
		"edit",              // edit tunnel
		"1",                 // tunnel index
		"web",               // name
		"web.example.com",   // hostname
		"127.0.0.1:8080",    // local
		"public",            // access
		"add",               // add tunnel
		"admin",             // name
		"",                  // hostname
		"127.0.0.1:9090",    // local
		"token_protected",   // access
		"BLOOP_ADMIN_TOKEN", // token env
		"done",              // finish
	}, "\n") + "\n")

	err := Run(&stdout, &stderr, Options{ConfigPath: "-", OutputMode: OutputYAML, Stdin: input})
	if err != nil {
		t.Fatalf("Run returned error: %v", err)
	}
	got := stdout.String()
	for _, needle := range []string{
		"name: web",
		"hostname: web.example.com",
		"local_addr: 127.0.0.1:8080",
		"name: admin",
		"token_env: BLOOP_ADMIN_TOKEN",
	} {
		if !strings.Contains(got, needle) {
			t.Fatalf("output missing %q:\n%s", needle, got)
		}
	}
}

func TestRunInteractiveDockerDiscoveryRequiresConfirmation(t *testing.T) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	input := strings.NewReader(strings.Join([]string{
		"", "", "", "", "", "", // defaults
		"y",    // run discovery
		"y",    // add first candidate
		"n",    // skip second
		"done", // finish tunnel editing
	}, "\n") + "\n")

	err := Run(&stdout, &stderr, Options{
		ConfigPath:     "-",
		OutputMode:     OutputYAML,
		Stdin:          input,
		DiscoverDocker: true,
		Discoverer: fakeDiscoverer{candidates: []dockerdiscover.Candidate{
			{ContainerName: "frontend", Image: "app:latest", LocalAddr: "frontend:3000", SuggestedName: "frontend-3000"},
			{ContainerName: "api", Image: "api:latest", LocalAddr: "api:8080", SuggestedName: "api-8080"},
		}},
	})
	if err != nil {
		t.Fatalf("Run returned error: %v", err)
	}
	got := stdout.String()
	if !strings.Contains(got, "name: frontend-3000") {
		t.Fatalf("expected discovered tunnel in output:\n%s", got)
	}
	if strings.Contains(got, "name: api-8080") {
		t.Fatalf("did not expect skipped discovered tunnel in output:\n%s", got)
	}
}

func TestRunInteractiveDockerDiscoveryGracefulFallback(t *testing.T) {
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
		t.Fatalf("expected graceful discovery warning, got %q", stderr.String())
	}
}
