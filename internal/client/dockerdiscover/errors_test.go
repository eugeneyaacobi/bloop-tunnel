package dockerdiscover

import (
	"context"
	"strings"
	"testing"
)

func TestDiscoverMissingSocketReturnsHelpfulError(t *testing.T) {
	client := NewClient(t.TempDir() + "/missing.sock")
	_, err := client.Discover(context.Background())
	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.Error(), "docker socket unavailable") {
		t.Fatalf("unexpected error: %v", err)
	}
}
