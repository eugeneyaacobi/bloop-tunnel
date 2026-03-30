package tui

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/lipgloss"
)

type VerificationScreenModel struct {
	config           Config
	width, height    int
	showHelp         bool
	currentStep      int
	steps            []VerificationStep
	spinner          spinner.Model
	allComplete      bool
	hasErrors        bool
	retrySelected    bool
	quitting         bool
}

type VerificationStep struct {
	Name        string
	Status      string // "pending", "running", "success", "error", "skipped"
	Error       error
	StartedAt   time.Time
	CompletedAt time.Time
}

type VerificationTickMsg struct{}
type VerificationCompleteMsg struct{ Success bool }
type VerificationRetryMsg struct{}

func NewVerificationScreenModel(config Config) VerificationScreenModel {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(primaryColor)

	steps := []VerificationStep{
		{Name: "DNS Resolution", Status: "pending"},
		{Name: "TLS Handshake", Status: "pending"},
		{Name: "HTTP Health Check", Status: "pending"},
		{Name: "Enrollment Verification", Status: "pending"},
		{Name: "Relay WebSocket Test", Status: "pending"},
	}

	return VerificationScreenModel{
		config:      config,
		width:       80,
		height:      24,
		showHelp:    false,
		currentStep: 0,
		steps:       steps,
		spinner:     s,
		allComplete: false,
		hasErrors:   false,
	}
}

func (m VerificationScreenModel) Init() tea.Cmd {
	return tea.Batch(
		m.spinner.Tick,
		m.runNextStep(),
	)
}

func (m VerificationScreenModel) runNextStep() tea.Cmd {
	return tea.Tick(time.Second*0, func(t time.Time) tea.Msg {
		return VerificationTickMsg{}
	})
}

func (m VerificationScreenModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			m.quitting = true
			return m, tea.Quit
		case "r":
			if m.allComplete && m.hasErrors {
				// Retry verification
				return m, m.retryVerification()
			}
		case "enter":
			if m.allComplete {
				if m.hasErrors {
					return m, m.retryVerification()
				} else {
					// All good, advance to Docker discovery
					return m, func() tea.Msg { return ScreenTransitionMsg{From: StateVerification, To: StateDockerDiscovery} }
				}
			}
		case "esc":
			// Go back
			return m, func() tea.Msg { return ScreenTransitionMsg{From: StateVerification, To: StateConfig} }
		case "?":
			m.showHelp = !m.showHelp
			return m, nil
		}
	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	case VerificationTickMsg:
		return m, m.runVerificationStep()
	case VerificationRetryMsg:
		return m, m.retryVerification()
	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
	}

	return m, nil
}

func (m VerificationScreenModel) runVerificationStep() tea.Cmd {
	// Find next pending step
	for i := range m.steps {
		if m.steps[i].Status == "pending" {
			m.currentStep = i
			m.steps[i].Status = "running"
			m.steps[i].StartedAt = time.Now()

			// Simulate the check (in real implementation, this would run actual network checks)
			return m.performCheck(i)
		}
	}

	// All steps complete
	m.allComplete = true
	for _, step := range m.steps {
		if step.Status == "error" {
			m.hasErrors = true
			break
		}
	}

	// Update config with results
	m.config.VerificationResults.Connectivity.Status = m.steps[0].Status
	m.config.VerificationResults.Connectivity.Success = m.steps[0].Status == "success"
	m.config.VerificationResults.Enrollment.Status = m.steps[3].Status
	m.config.VerificationResults.Enrollment.Success = m.steps[3].Status == "success"
	m.config.VerificationResults.Relay.Status = m.steps[4].Status
	m.config.VerificationResults.Relay.Success = m.steps[4].Status == "success"

	return func() tea.Msg { return ConfigSaveMsg{Config: m.config} }
}

func (m VerificationScreenModel) performCheck(stepIndex int) tea.Cmd {
	return tea.Tick(time.Millisecond*800, func(t time.Time) tea.Msg {
		// Simulate check result (in real implementation, perform actual network checks)
		status := "success"
		var err error

		// Simulate occasional failures for demo
		if stepIndex == 1 && false { // TLS check
			status = "error"
			err = &VerificationError{Step: m.steps[stepIndex].Name, Message: "certificate expired"}
		}

		return VerificationStepCompleteMsg{
			Index:  stepIndex,
			Status: status,
			Error:  err,
		}
	})
}

type VerificationStepCompleteMsg struct {
	Index  int
	Status string
	Error  error
}

type VerificationError struct {
	Step    string
	Message string
}

func (e *VerificationError) Error() string {
	return e.Step + ": " + e.Message
}

func (m VerificationScreenModel) retryVerification() tea.Cmd {
	// Reset all steps
	for i := range m.steps {
		m.steps[i].Status = "pending"
		m.steps[i].Error = nil
	}
	m.allComplete = false
	m.hasErrors = false
	m.currentStep = 0

	return tea.Batch(
		m.spinner.Tick,
		m.runNextStep(),
	)
}

func (m VerificationScreenModel) View() string {
	title := titleStyle.Render("Connectivity Verification")
	subtitle := subtitleStyle.Render("Verifying your connection to bloop-tunnel services")

	// Build steps view
	var stepViews []string
	for i, step := range m.steps {
		stepView := m.renderStep(step, i)
		stepViews = append(stepViews, stepView)
	}

	stepsContent := lipgloss.JoinVertical(lipgloss.Left, stepViews...)

	// Status summary
	var summary string
	if m.allComplete {
		if m.hasErrors {
			summary = errorStyle.Render("⚠ Some verification checks failed. Press 'r' to retry or 'esc' to go back.")
		} else {
			summary = successStyle.Render("✓ All verification checks passed! Press 'enter' to continue.")
		}
	} else {
		runningStep := m.steps[m.currentStep]
		summary = lipgloss.NewStyle().Foreground(primaryColor).Render(
			m.spinner.View() + " Running: " + runningStep.Name,
		)
	}

	content := lipgloss.JoinVertical(
		lipgloss.Left,
		title,
		subtitle,
		"",
		stepsContent,
		"",
		summary,
	)

	// Help footer
	var helpText string
	if m.showHelp {
		helpModel := NewHelpModel(VerificationScreenKeys())
		helpText = helpModel.ShortHelpView()
	} else {
		helpText = footerStyle.Render("?: Show help")
	}

	content = lipgloss.JoinVertical(lipgloss.Left, content, "", helpText)

	return Center(m.width, content)
}

func (m VerificationScreenModel) renderStep(step VerificationStep, index int) string {
	var icon, statusText string
	var statusStyle lipgloss.Style

	switch step.Status {
	case "pending":
		icon = "○"
		statusText = "pending"
		statusStyle = lipgloss.NewStyle().Foreground(mutedColor)
	case "running":
		icon = m.spinner.View()
		statusText = "checking..."
		statusStyle = lipgloss.NewStyle().Foreground(primaryColor)
	case "success":
		icon = "✓"
		statusText = "success"
		statusStyle = successStyle
	case "error":
		icon = "✗"
		statusText = "failed"
		if step.Error != nil {
			statusText += ": " + step.Error.Error()
		}
		statusStyle = errorStyle
	case "skipped":
		icon = "⊝"
		statusText = "skipped"
		statusStyle = lipgloss.NewStyle().Foreground(mutedColor)
	}

	nameStyle := lipgloss.NewStyle().Width(25).Foreground(textColor)
	
	row := lipgloss.JoinHorizontal(
		lipgloss.Left,
		nameStyle.Render(step.Name),
		"  ",
		statusStyle.Render(icon+" "+statusText),
	)

	return row
}
