package models

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type StatusModel struct {
	loading  bool
	message  string
	frame    int   // For loading animation
	width     int
}

type StatusLoadingMsg struct{ Message string }
type StatusSuccessMsg struct{ Message string }
type StatusErrorMsg struct{ Message string }

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
	case StatusSuccessMsg:
		m.loading = false
		m.message = msg.Message
	case StatusErrorMsg:
		m.loading = false
		m.message = msg.Message
	case tea.WindowSizeMsg:
		m.width = msg.Width
	case tea.KeyMsg:
		// Animate loading frame on any keypress
		if m.loading {
			m.frame = (m.frame + 1) % 4
		}
	}

	return m, nil
}

func (m *StatusModel) View() string {
	if m.loading {
		// Simple loading animation
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
