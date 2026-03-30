package tui

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

var (
	// Colors — adaptive for light/dark terminals
	primaryColor = lipgloss.AdaptiveColor{Light: "#0055D4", Dark: "#007AFF"}
	successColor = lipgloss.AdaptiveColor{Light: "#228B22", Dark: "#34C759"}
	errorColor   = lipgloss.AdaptiveColor{Light: "#CC0000", Dark: "#FF3B30"}
	warningColor = lipgloss.AdaptiveColor{Light: "#B8860B", Dark: "#FF9500"}
	mutedColor   = lipgloss.AdaptiveColor{Light: "#666666", Dark: "#8E8E93"}
	textColor    = lipgloss.AdaptiveColor{Light: "#1A1A1A", Dark: "#FFFFFF"}
	bgColor      = lipgloss.AdaptiveColor{Light: "#F5F5F5", Dark: "#1E1E1E"}

	// Component Styles
	titleStyle        = lipgloss.NewStyle().Foreground(primaryColor).Bold(true).MarginTop(1).MarginBottom(1)
	subtitleStyle     = lipgloss.NewStyle().Foreground(mutedColor).Italic(true).MarginBottom(2)
	inputStyle        = lipgloss.NewStyle().Foreground(textColor)
	inputFocusedStyle = lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).Foreground(textColor)
	buttonStyle       = lipgloss.NewStyle().Padding(0, 2)
	buttonFocusedStyle = lipgloss.NewStyle().Padding(0, 2).Foreground(textColor).Background(primaryColor).Bold(true)
	successStyle      = lipgloss.NewStyle().Foreground(successColor)
	errorStyle        = lipgloss.NewStyle().Foreground(errorColor)
	warningStyle      = lipgloss.NewStyle().Foreground(warningColor)
	listItemStyle     = lipgloss.NewStyle().Padding(0, 1)
	listFocusedStyle  = lipgloss.NewStyle().Padding(0, 1).Background(mutedColor).Foreground(textColor)
	footerStyle       = lipgloss.NewStyle().Foreground(mutedColor).MarginTop(2)
	borderStyle       = lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).BorderForeground(primaryColor)
	boxStyle          = lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).BorderForeground(primaryColor).Padding(1, 2)
	descStyle         = lipgloss.NewStyle().Foreground(mutedColor).Italic(true).MarginLeft(4)

	// Shared box styles for sub-components
	componentBoxStyle    = lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).BorderForeground(primaryColor).Padding(0, 1).Width(52)
	componentBoxFocused  = lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).BorderForeground(primaryColor).Foreground(primaryColor).Bold(true).Padding(0, 1).Width(52)
	dialogBoxStyle       = lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).BorderForeground(primaryColor).Padding(1, 2).Width(52).Align(lipgloss.Center)
	optionStyle          = lipgloss.NewStyle().Padding(0, 2).MarginRight(2).Width(25)
	optionSelectedStyle  = lipgloss.NewStyle().Padding(0, 2).MarginRight(2).Width(25).Background(primaryColor).Foreground(textColor).Bold(true)
	labelStyle           = lipgloss.NewStyle().Foreground(primaryColor).Bold(true)
	placeholderStyle     = lipgloss.NewStyle().Foreground(mutedColor).Italic(true)
	valueStyle           = lipgloss.NewStyle().Foreground(textColor)
	fieldErrorStyle      = lipgloss.NewStyle().Foreground(errorColor)

	// Radio/dot indicators
	dotSelected   = "● "
	dotUnselected = "○ "
	dotNone       = "  "
)

// Layout helpers

// Center content horizontally
func Center(width int, content string) string {
	return lipgloss.Place(width, 0, lipgloss.Center, lipgloss.Top, content)
}

// Horizontal separator
func HSeparator(width int) string {
	return lipgloss.NewStyle().Foreground(mutedColor).Width(width).Render(strings.Repeat("─", width))
}

// Section header
func SectionHeader(title string) string {
	return lipgloss.NewStyle().
		Foreground(primaryColor).
		Bold(true).
		MarginTop(1).
		MarginBottom(1).
		Render(title)
}

// Label with colon
func Label(label string) string {
	return labelStyle.Render(label + ":")
}

// Value display
func Value(value string) string {
	return valueStyle.Render(value)
}

// Password placeholder
func PasswordMask(length int) string {
	if length == 0 {
		return "(empty)"
	}
	return strings.Repeat("•", length)
}

// Key shortcut
func KeyShortcut(keys string) string {
	return lipgloss.NewStyle().
		Foreground(primaryColor).
		Bold(true).
		Render(keys)
}

// Help text
func HelpText(text string) string {
	return lipgloss.NewStyle().
		Foreground(mutedColor).
		Italic(true).
		Render(text)
}

// Success message
func SuccessMessage(text string) string {
	return successStyle.Render("✓ " + text)
}

// Error message
func ErrorMessage(text string) string {
	return errorStyle.Render("✗ " + text)
}

// Warning message
func WarningMessage(text string) string {
	return warningStyle.Render("⚠ " + text)
}
