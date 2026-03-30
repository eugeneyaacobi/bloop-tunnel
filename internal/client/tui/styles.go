package tui

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

var (
	// Colors
	primaryColor    = lipgloss.Color("#007AFF")  // Blue
	successColor    = lipgloss.Color("#34C759")  // Green
	errorColor      = lipgloss.Color("#FF3B30")  // Red
	warningColor   = lipgloss.Color("#FF9500")  // Orange
	mutedColor     = lipgloss.Color("#8E8E93")  // Gray
	textColor      = lipgloss.Color("#FFFFFF")  // White
	bgColor        = lipgloss.Color("#1E1E1E")  // Dark background

	// Component Styles
	titleStyle       = lipgloss.NewStyle().Foreground(primaryColor).Bold(true).MarginTop(1).MarginBottom(1)
	subtitleStyle    = lipgloss.NewStyle().Foreground(mutedColor).Italic(true).MarginBottom(2)
	inputStyle       = lipgloss.NewStyle().Foreground(textColor)
	inputFocusedStyle = lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).Foreground(textColor)
	buttonStyle      = lipgloss.NewStyle().Padding(0, 2)
	buttonFocusedStyle = buttonStyle.Copy().Foreground(textColor).Background(primaryColor).Bold(true)
	successStyle     = lipgloss.NewStyle().Foreground(successColor)
	errorStyle       = lipgloss.NewStyle().Foreground(errorColor)
	listItemStyle    = lipgloss.NewStyle().Padding(0, 1)
	listFocusedStyle = listItemStyle.Copy().Background(mutedColor).Foreground(textColor)
	footerStyle       = lipgloss.NewStyle().Foreground(mutedColor).MarginTop(2)
	borderStyle       = lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).BorderForeground(primaryColor)
	boxStyle        = lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).BorderForeground(primaryColor).Padding(1, 2)
	descStyle       = lipgloss.NewStyle().Foreground(mutedColor).Italic(true).MarginLeft(4)
)

// Layout helpers

// Center content horizontally
func Center(width int, content string) string {
	return lipgloss.Place(width, 0, lipgloss.Center, lipgloss.Top, content)
}

// Vertical separator
func VSeparator() string {
	return lipgloss.NewStyle().Foreground(mutedColor).Render("─")
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

// Section content
func SectionContent(content string) string {
	return lipgloss.NewStyle().
		Foreground(textColor).
		MarginLeft(2).
		Render(content)
}

// Section error
func SectionError(err string) string {
	return lipgloss.NewStyle().
		Foreground(errorColor).
		MarginTop(1).
		MarginLeft(2).
		Render("  └─ " + err)
}

// Label with colon
func Label(label string) string {
	return lipgloss.NewStyle().
		Foreground(primaryColor).
		Bold(true).
		Render(label + ":")
}

// Value display
func Value(value string) string {
	return lipgloss.NewStyle().
		Foreground(textColor).
		Render(value)
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
	return lipgloss.NewStyle().
		Foreground(successColor).
		Render("✓ " + text)
}

// Error message
func ErrorMessage(text string) string {
	return lipgloss.NewStyle().
		Foreground(errorColor).
		Render("✗ " + text)
}

// Warning message
func WarningMessage(text string) string {
	return lipgloss.NewStyle().
		Foreground(warningColor).
		Render("⚠ " + text)
}
