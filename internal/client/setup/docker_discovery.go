package setup

import (
	"bufio"
	"context"
	"fmt"
	"io"

	"bloop-tunnel/internal/client/dockerdiscover"
	"bloop-tunnel/internal/config"
)

func applyDockerDiscovery(ctx context.Context, reader *bufio.Reader, out io.Writer, cfg *config.ClientConfig, opts Options) error {
	discoverer := opts.Discoverer
	if discoverer == nil {
		discoverer = dockerdiscover.NewClient("")
	}
	candidates, err := discoverer.Discover(ctx)
	if err != nil {
		return err
	}
	if len(candidates) == 0 {
		fmt.Fprintln(out, "No eligible Docker services found. Continuing with manual tunnel entry.")
		return nil
	}
	fmt.Fprintln(out, "Docker discovery candidates:")
	for i, candidate := range candidates {
		fmt.Fprintf(out, "  [%d] %s (%s) -> %s\n", i+1, candidate.ContainerName, candidate.Image, candidate.LocalAddr)
		if !confirm(reader, out, fmt.Sprintf("Add tunnel for %s?", candidate.LocalAddr), false) {
			continue
		}
		cfg.Tunnels = append(cfg.Tunnels, config.TunnelConfig{
			Name:      candidate.SuggestedName,
			LocalAddr: candidate.LocalAddr,
			Access:    "public",
		})
	}
	return nil
}
