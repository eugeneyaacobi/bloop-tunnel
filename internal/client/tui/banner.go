package tui

import (
	"fmt"
	"strings"
)

// bloopArt in banner font — chunky, immediately readable.
var bloopArt = []string{
	"#####  #       ####   ####  #####  ",
	"#    # #      #    # #    # #    # ",
	"#####  #      #    # #    # #    # ",
	"#    # #      #    # #    # #####  ",
	"#    # #      #    # #    # #      ",
	"#####  ######  ####   ####  #      ",
}

// tunnelArt in banner font.
var tunnelArt = []string{
	"##### #    # #    # #    # ###### #       ",
	"  #   #    # ##   # ##   # #      #       ",
	"  #   #    # # #  # # #  # #####  #       ",
	"  #   #    # #  # # #  # # #      #       ",
	"  #   #    # #   ## #   ## #      #       ",
	"  #    ####  #    # #    # ###### ######  ",
}

// gradientColors defines the indigo→cyan gradient (#6366f1 → #06b6d4).
var gradientColors = [][3]int{
	{99, 102, 241},
	{88, 110, 237},
	{77, 118, 234},
	{66, 126, 231},
	{55, 134, 228},
	{44, 142, 224},
	{33, 150, 221},
	{22, 158, 218},
	{11, 166, 215},
	{6, 174, 213},
	{6, 182, 212},
}

// applyGradient applies a horizontal ANSI 24-bit color gradient to a line of text.
func applyGradient(line string) string {
	var sb strings.Builder
	width := len(line)
	for i, ch := range line {
		if ch == ' ' {
			sb.WriteByte(' ')
			continue
		}
		t := float64(i) / float64(width)
		idx := t * float64(len(gradientColors)-1)
		lo := int(idx)
		hi := lo + 1
		if hi >= len(gradientColors) {
			hi = len(gradientColors) - 1
		}
		frac := idx - float64(lo)
		r := int(float64(gradientColors[lo][0])*(1-frac) + float64(gradientColors[hi][0])*frac)
		g := int(float64(gradientColors[lo][1])*(1-frac) + float64(gradientColors[hi][1])*frac)
		b := int(float64(gradientColors[lo][2])*(1-frac) + float64(gradientColors[hi][2])*frac)
		fmt.Fprintf(&sb, "\x1b[38;2;%d;%d;%dm%c\x1b[0m", r, g, b, ch)
	}
	return sb.String()
}

func init() {
	var sb strings.Builder
	sb.WriteString("\n")
	for _, line := range bloopArt {
		sb.WriteString("  ")
		sb.WriteString(applyGradient(line))
		sb.WriteString("\n")
	}
	sb.WriteString("\n")
	for _, line := range tunnelArt {
		sb.WriteString("  ")
		sb.WriteString(applyGradient(line))
		sb.WriteString("\n")
	}
	sb.WriteString("\n")
	sb.WriteString("  \x1b[38;2;99;102;241mSecure tunnels for localhost\x1b[0m\n")
	sb.WriteString("  \x1b[2mZero config. Just run it.\x1b[0m\n")
	sb.WriteString("\n")
	WelcomeBanner = sb.String()
}

// WelcomeBanner is the ASCII art banner displayed at the start of the TUI setup wizard.
var WelcomeBanner string
