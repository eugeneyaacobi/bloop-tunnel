package models

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// ReviewScreenConfirmMsg is sent when user confirms the configuration
type ReviewScreenConfirmMsg struct{}

// ReviewScreenBackMsg is sent when user wants to go back
type ReviewScreenBackMsg struct{}

// ReviewScreenEditMsg is sent when user wants to edit a section
type ReviewScreenEditMsg struct {
	Section string
}

// ReviewScreenData represents the configuration to review
type ReviewScreenData struct {
	Global      GlobalConfigReview
	Tunnels     []TunnelConfigReview
	Output      string
	OutputPath  string
}

type GlobalConfigReview struct {
	ControlPlaneURL string
	RelayURL        string
	AuthTokenEnv    string
}

type TunnelConfigReview struct {
	Name      string
	Hostname  string
	LocalAddr string
	Access    string
	AuthDetails string
}

// ReviewScreenModel for reviewing configuration before output
type ReviewScreenModel struct {
	data       ReviewScreenData
	width      int
	height     int
	scrollPos  int
	visibleHeight int
}

// ReviewScreenOpts for constructing the model
type ReviewScreenOpts struct {
	Data ReviewScreenData
}

func NewReviewScreenModel(opts ReviewScreenOpts) ReviewScreenModel {
	return ReviewScreenModel{
		data:           opts.Data,
		width:          80,
		height:         24,
		scrollPos:      0,
		visibleHeight:  15,
	}
}

func (m ReviewScreenModel) Init() tea.Cmd {
	return nil
}

func (m ReviewScreenModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "b", "esc":
			return m, func() tea.Msg { return ReviewScreenBackMsg{} }
		case "enter":
			return m, func() tea.Msg { return ReviewScreenConfirmMsg{} }
		case "e":
			// Edit global config
			return m, func() tea.Msg { return ReviewScreenEditMsg{Section: "global"} }
		case "up", "k":
			if m.scrollPos > 0 {
				m.scrollPos--
			}
		case "down", "j":
			maxScroll := m.totalLines() - m.visibleHeight
			if m.scrollPos < maxScroll {
				m.scrollPos++
			}
		case "pageup", "ctrl+b":
			m.scrollPos = max(0, m.scrollPos-m.visibleHeight/2)
		case "pagedown", "ctrl+f":
			maxScroll := m.totalLines() - m.visibleHeight
			m.scrollPos = min(maxScroll, m.scrollPos+m.visibleHeight/2)
		case "home", "g":
			m.scrollPos = 0
		case "end", "G":
			maxScroll := m.totalLines() - m.visibleHeight
			m.scrollPos = max(0, maxScroll)
		}

	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
		m.visibleHeight = m.height - 8
		if m.visibleHeight < 5 {
			m.visibleHeight = 5
		}
	}

	return m, nil
}

func (m ReviewScreenModel) totalLines() int {
	// Estimate total lines for scrolling
	baseLines := 10 // Header and global config
	tunnelLines := len(m.data.Tunnels) * 4 // Each tunnel takes ~4 lines
	return baseLines + tunnelLines
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func (m ReviewScreenModel) View() string {
	titleStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#007AFF")).
		Bold(true).
		MarginTop(1).
		MarginBottom(1)

	headerStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#8E8E93")).
		Italic(true).
		MarginBottom(1)

	sectionStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#007AFF")).
		Bold(true).
		MarginTop(1).
		MarginBottom(0)

	labelStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#8E8E93")).
		Width(20)

	valueStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FFFFFF"))

	tunnelStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#32D74B")).
		Bold(true).

		MarginLeft(2)

	footerStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#8E8E93")).
		MarginTop(1)

	var sections []string

	// Header
	sections = append(sections, titleStyle.Render("Review Configuration"))
	sections = append(sections, "")
	sections = append(sections, headerStyle.Render("Review your settings before generating output:"))
	sections = append(sections, "")

	// Global Configuration
	sections = append(sections, sectionStyle.Render("Global Configuration"))
	sections = append(sections, labelStyle.Render("Control Plane:") + valueStyle.Render(m.data.Global.ControlPlaneURL))
	sections = append(sections, labelStyle.Render("Relay URL:") + valueStyle.Render(m.data.Global.RelayURL))
	sections = append(sections, labelStyle.Render("Auth Token Env:") + valueStyle.Render(m.data.Global.AuthTokenEnv))
	sections = append(sections, "")

	// Tunnels
	sections = append(sections, sectionStyle.Render(fmt.Sprintf("Tunnels (%d)", len(m.data.Tunnels))))
	if len(m.data.Tunnels) == 0 {
		sections = append(sections, labelStyle.Render("")+valueStyle.Foreground(lipgloss.Color("#8E8E93")).Render("No tunnels configured"))
	} else {
		for _, t := range m.data.Tunnels {
			sections = append(sections, tunnelStyle.Render("● "+t.Name))
			sections = append(sections, labelStyle.Render("  Local:") + valueStyle.Render(t.LocalAddr))
			if t.Hostname != "" {
				sections = append(sections, labelStyle.Render("  Hostname:") + valueStyle.Render(t.Hostname))
			}
			sections = append(sections, labelStyle.Render("  Access:") + valueStyle.Render(t.Access))
			if t.AuthDetails != "" {
				sections = append(sections, labelStyle.Render("  Auth:") + valueStyle.Render(t.AuthDetails))
			}
		}
	}
	sections = append(sections, "")

	// Output
	sections = append(sections, sectionStyle.Render("Output"))
	sections = append(sections, labelStyle.Render("Format:") + valueStyle.Render(m.data.Output))
	sections = append(sections, labelStyle.Render("Path:") + valueStyle.Render(m.data.OutputPath))

	// Scroll content
	allLines := sections
	visibleEnd := min(m.scrollPos+m.visibleHeight, len(allLines))
	visibleLines := allLines[m.scrollPos:visibleEnd]

	// Add scroll indicators
	if m.scrollPos > 0 {
		visibleLines = append([]string{"▲"}, visibleLines...)
	}
	if visibleEnd < len(allLines) {
		visibleLines = append(visibleLines, "▼")
	}

	return lipgloss.JoinVertical(
		lipgloss.Left,
		lipgloss.JoinVertical(lipgloss.Left, visibleLines...),
		"",
		footerStyle.Render("↑/↓: Scroll • E: Edit config • Enter: Confirm • B: Back • Q: Quit"),
	)
}
