package tui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/bubbles/viewport"
	"github.com/charmbracelet/lipgloss"
)

type ReviewScreenModel struct {
	config      Config
	viewport    viewport.Model
	width       int
	height      int
	showHelp    bool
	outputMode  OutputMode
}

func NewReviewScreenModel(config Config) ReviewScreenModel {
	vp := viewport.New(80, 20)
	vp.Style = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(primaryColor).
		Padding(1, 1)

	return ReviewScreenModel{
		config:     config,
		viewport:   vp,
		width:      80,
		height:     24,
		showHelp:   false,
		outputMode: config.OutputMode,
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
		case "enter":
			// Confirm and generate output
			return m, func() tea.Msg { return ScreenTransitionMsg{From: StateReview, To: StateOutput} }
		case "esc", "b":
			// Go back to edit
			return m, func() tea.Msg { return ScreenTransitionMsg{From: StateReview, To: StateConfig} }
		case "1", "2", "3":
			// Change output mode
			switch msg.String() {
			case "1":
				m.outputMode = OutputYAML
			case "2":
				m.outputMode = OutputEnvFile
			case "3":
				m.outputMode = OutputComposeBlock
			}
			m.config.OutputMode = m.outputMode
		case "?":
			m.showHelp = !m.showHelp
			return m, nil
		}
	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
		m.viewport.Width = m.width - 4
		m.viewport.Height = m.height - 12
		m.viewport.GotoTop()
	}

	var cmd tea.Cmd
	m.viewport, cmd = m.viewport.Update(msg)
	return m, cmd
}

func (m ReviewScreenModel) View() string {
	title := titleStyle.Render("Configuration Review")
	subtitle := subtitleStyle.Render("Review your configuration before generating output")

	// Build review content
	content := m.buildReviewContent()
	m.viewport.SetContent(content)

	// Output mode selector
	modeSelector := m.renderOutputModeSelector()

	// Action buttons
	backButton := buttonStyle.Render("← Back (Esc)")
	confirmButton := buttonFocusedStyle.Render("Confirm (Enter)")

	buttonRow := lipgloss.JoinHorizontal(lipgloss.Left, backButton, "  ", confirmButton)

	// Main content layout
	mainContent := lipgloss.JoinVertical(
		lipgloss.Left,
		title,
		subtitle,
		"",
		modeSelector,
		"",
		m.viewport.View(),
		"",
		buttonRow,
	)

	// Help footer
	var helpText string
	if m.showHelp {
		helpText = HelpText("↑/↓: Scroll • 1/2/3: Output mode • Enter: Confirm • Esc: Back")
	} else {
		helpText = footerStyle.Render("?: Show help")
	}

	fullContent := lipgloss.JoinVertical(lipgloss.Left, mainContent, "", helpText)

	return lipgloss.NewStyle().
		Width(m.width).
		Align(lipgloss.Center).
		Render(fullContent)
}

func (m ReviewScreenModel) buildReviewContent() string {
	var sections []string

	// Connection Settings section
	sections = append(sections, SectionHeader("Connection Settings"))
	sections = append(sections, Label("Control Plane")+": "+Value(m.config.ControlPlaneURL))
	sections = append(sections, Label("Relay URL")+": "+Value(m.config.RelayURL))
	sections = append(sections, Label("Auth Token Env")+": "+Value(m.config.AuthTokenEnv))
	if m.config.EnrollmentToken != "" {
		sections = append(sections, Label("Enrollment Token")+": "+Value("******"))
	}
	sections = append(sections, Label("Reconnect")+": "+Value(fmt.Sprintf("%d-%dms",
		m.config.Reconnect.InitialDelayMs, m.config.Reconnect.MaxDelayMs)))

	// Verification Results section
	sections = append(sections, "")
	sections = append(sections, SectionHeader("Verification Results"))

	if m.config.VerificationResults.Connectivity.Status != "" {
		status := m.config.VerificationResults.Connectivity.Status
		if status == "success" {
			sections = append(sections, SuccessMessage("Connectivity: OK"))
		} else if status == "error" {
			sections = append(sections, ErrorMessage("Connectivity: Failed"))
		} else {
			sections = append(sections, WarningMessage("Connectivity: "+status))
		}
	}

	if m.config.VerificationResults.Enrollment.Status != "" {
		status := m.config.VerificationResults.Enrollment.Status
		if status == "success" {
			sections = append(sections, SuccessMessage("Enrollment: OK"))
			if m.config.VerificationResults.Enrollment.InstallationID != "" {
				sections = append(sections, Value("  Installation ID: "+m.config.VerificationResults.Enrollment.InstallationID))
			}
		} else if status == "error" {
			sections = append(sections, ErrorMessage("Enrollment: Failed"))
		} else {
			sections = append(sections, WarningMessage("Enrollment: "+status))
		}
	}

	if m.config.VerificationResults.Relay.Status != "" {
		status := m.config.VerificationResults.Relay.Status
		if status == "success" {
			sections = append(sections, SuccessMessage("Relay Connection: OK"))
		} else if status == "error" {
			sections = append(sections, ErrorMessage("Relay Connection: Failed"))
		} else {
			sections = append(sections, WarningMessage("Relay Connection: "+status))
		}
	}

	// Tunnels section
	sections = append(sections, "")
	sections = append(sections, SectionHeader("Configured Tunnels"))

	if len(m.config.Tunnels) == 0 {
		sections = append(sections, HelpText("No tunnels configured"))
	} else {
		for i, tunnel := range m.config.Tunnels {
			sections = append(sections, "")
			sections = append(sections, fmt.Sprintf("%d. %s", i+1, SuccessMessage(tunnel.Name)))
			sections = append(sections, "   "+Label("Local")+": "+Value(fmt.Sprintf("%s:%d", tunnel.LocalIP, tunnel.LocalPort)))

			if tunnel.Hostname != "" {
				sections = append(sections, "   "+Label("Hostname")+": "+Value(tunnel.Hostname))
			}

			accessInfo := ""
			switch tunnel.Access {
			case "public":
				accessInfo = "🌐 Public"
			case "basic_auth":
				accessInfo = fmt.Sprintf("🔒 Basic Auth (%s)", tunnel.BasicAuth.Username)
			case "token_protected":
				accessInfo = fmt.Sprintf("🔑 Token (%s)", tunnel.TokenEnv)
			}
			sections = append(sections, "   "+Label("Access")+": "+Value(accessInfo))
		}
	}

	// Output section
	sections = append(sections, "")
	sections = append(sections, SectionHeader("Output"))
	sections = append(sections, Label("Format")+": "+Value(string(m.outputMode)))
	sections = append(sections, Label("Path")+": "+Value(m.config.OutputPath))

	return strings.Join(sections, "\n")
}

func (m ReviewScreenModel) renderOutputModeSelector() string {
	modes := []struct {
		Key  string
		Name string
		Mode OutputMode
	}{
		{"1", "YAML", OutputYAML},
		{"2", "Env File", OutputEnvFile},
		{"3", "Docker Compose", OutputComposeBlock},
	}

	var buttons []string
	for _, mode := range modes {
		btnStyle := buttonStyle
		if m.outputMode == mode.Mode {
			btnStyle = buttonFocusedStyle
		}
		buttons = append(buttons, btnStyle.Render(fmt.Sprintf("%s: %s", mode.Key, mode.Name)))
	}

	return lipgloss.JoinHorizontal(lipgloss.Left, buttons...)
}

func (m ReviewScreenModel) WithConfig(config Config) ReviewScreenModel {
	m.config = config
	m.outputMode = config.OutputMode
	return m
}
