package models

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type StatusModel struct {
	loading bool
	message string
	frame   int
	width   int
}

type StatusLoadingMsg struct{ Message string }
type StatusSuccessMsg struct{ Message string }
type StatusErrorMsg struct{ Message string }

type spinnerTickMsg struct{}

func NewStatusModel() *StatusModel {
	return &StatusModel{
		loading: false,
		message: "",
		frame:   0,
		width:   60,
	}
}

func (m *StatusModel) Init() tea.Cmd {
	return nil
}

func (m *StatusModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case StatusLoadingMsg:
		m.loading = true
		m.message = msg.Message
		m.frame = 0
		return m, spinTickCmd()
	case StatusSuccessMsg:
		m.loading = false
		m.message = msg.Message
	case StatusErrorMsg:
		m.loading = false
		m.message = msg.Message
	case spinnerTickMsg:
		if m.loading {
			m.frame++
			return m, spinTickCmd()
		}
	case tea.WindowSizeMsg:
		m.width = msg.Width
	}

	return m, nil
}

func spinTickCmd() tea.Cmd {
	return tea.Tick(80*time.Millisecond, func(_ time.Time) tea.Msg {
		return spinnerTickMsg{}
	})
}

func (m *StatusModel) View() string {
	if m.loading {
		frames := []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}
		frame := frames[m.frame%len(frames)]
		spinnerText := lipgloss.NewStyle().
			Foreground(lipgloss.Color("#007AFF")).
			Render(frame + " " + m.message)
		return lipgloss.NewStyle().
			Width(m.width).
			Align(lipgloss.Center).
			Render(spinnerText)
	}

	statusStyle := lipgloss.NewStyle().
		Width(m.width).
		Align(lipgloss.Center)

	return statusStyle.Render(m.message)
}

func (m *StatusModel) WithLoading(message string) *StatusModel {
	m.loading = true
	m.message = message
	return m
}

func (m *StatusModel) WithSuccess(message string) *StatusModel {
	m.loading = false
	m.message = message
	return m
}

func (m *StatusModel) WithError(message string) *StatusModel {
	m.loading = false
	m.message = message
	return m
}

func (m *StatusModel) Clear() *StatusModel {
	m.loading = false
	m.message = ""
	return m
}
