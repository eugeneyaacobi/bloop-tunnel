package dockerdiscover

import "testing"

func TestLooksHTTPish(t *testing.T) {
	tests := []struct {
		port int
		want bool
	}{
		{80, true},
		{3000, true},
		{8080, true},
		{5432, false}, // postgres excluded even in known range
		{22, false},   // below 1024
		{10001, false},
		{443, false},
		{5000, true},
		{8888, true},
		{1024, true},
	}
	for _, tt := range tests {
		if got := looksHTTPish(tt.port); got != tt.want {
			t.Errorf("looksHTTPish(%d) = %v, want %v", tt.port, got, tt.want)
		}
	}
}

func TestExcluded(t *testing.T) {
	tests := []struct {
		name  string
		image string
		want  bool
	}{
		{"app", "node:18", false},
		{"frontend", "react:latest", false},
		{"api", "golang:1.24", false},
		{"postgres", "postgres:16", true},
		{"my-db", "mysql:8", true},
		{"cache", "redis:7", true},
		{"queue", "rabbitmq:3", true},
		{"storage", "minio/minio:latest", true},
		{"broker", "nats:2", true},
		{"myapp-postgres", "postgres:16", true},
	}
	for _, tt := range tests {
		if got := excluded(tt.name, tt.image); got != tt.want {
			t.Errorf("excluded(%q, %q) = %v, want %v", tt.name, tt.image, got, tt.want)
		}
	}
}
