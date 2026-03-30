package tui

import "github.com/charmbracelet/lipgloss"

// bannerColors defines the ANSI color palette for the bloop tunnel logo.
var bannerColors = []lipgloss.Color{
	"#6C5CE7", // purple
	"#0984E3", // blue
	"#00B894", // teal
	"#00CEC9", // cyan
	"#74B9FF", // light blue
}

// WelcomeBanner returns the ANSI-colored ASCII art banner for bloop-tunnel.
// Uses lipgloss for consistent terminal color support.
var WelcomeBanner = buildBanner()

func buildBanner() string {
	lines := []string{
		"",
		"  ██████╗ ██╗   ██╗ ██████╗     ████████╗███████╗███╗   ███╗",
		"  ██╔══██╗██║   ██║██╔════╝     ╚══██╔══╝██╔════╝████╗ ████║",
		"  ██████╔╝██║   ██║██║  ███╗       ██║   █████╗  ██╔████╔██║",
		"  ██╔══██╗██║   ██║██║   ██║       ██║   ██╔══╝  ██║╚██╔╝██║",
		"  ██████╔╝╚██████╔╝╚██████╔╝       ██║   ███████╗██║ ╚═╝ ██║",
		"  ╚═════╝  ╚═════╝  ╚═════╝        ╚═╝   ╚══════╝╚═╝     ╚═╝",
		"",
		"            ███╗   ██╗███████╗██╗  ██╗██╗   ██╗███████╗",
		"            ████╗  ██║██╔════╝╚██╗██╔╝██║   ██║██╔════╝",
		"            ██╔██╗ ██║█████╗   ╚███╔╝ ██║   ██║███████╗",
		"            ██║╚██╗██║██╔══╝   ██╔██╗ ██║   ██║╚════██║",
		"            ██║ ╚████║███████╗██╔╝ ██╗╚██████╔╝███████║",
		"            ╚═╝  ╚═══╝╚══════╝╚═╝  ╚═╝ ╚═════╝ ╚══════╝",
		"",
	}

	var styled string
	for i, line := range lines {
		if line == "" {
			styled += "\n"
			continue
		}
		colorIdx := i % len(bannerColors)
		style := lipgloss.NewStyle().Foreground(bannerColors[colorIdx]).Bold(true)
		styled += style.Render(line) + "\n"
	}

	tagline := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#8E8E93")).
		Italic(true).
		Render("         expose localhost · zero config · secure tunnels")

	styled += tagline + "\n"
	return styled
}
