package models

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/lipgloss"
)

type StatusModel struct {
	loading bool
	message string
	spinner spinner.Model
	width   int
}

type StatusLoadingMsg struct{ Message string }
type StatusSuccessMsg struct{ Message string }
type StatusErrorMsg struct{ Message string }

func NewStatusModel() *StatusModel {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("#007AFF"))
	return &StatusModel{
		loading: false,
		message: "",
		spinner: s,
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
		return m, m.spinner.Tick
	case StatusSuccessMsg:
		m.loading = false
		m.message = msg.Message
	case StatusErrorMsg:
		m.loading = false
		m.message = msg.Message
	case spinner.TickMsg:
		if m.loading {
			var cmd tea.Cmd
			m.spinner, cmd = m.spinner.Update(msg)
			return m, cmd
		}
	case tea.WindowSizeMsg:
		m.width = msg.Width
	}

	return m, nil
}

func (m *StatusModel) View() string {
	if m.loading {
		text := m.spinner.View() + " " + m.message
		return lipgloss.NewStyle().
			Width(m.width).
			Align(lipgloss.Center).
			Render(text)
	}

	return lipgloss.NewStyle().
		Width(m.width).
		Align(lipgloss.Center).
		Render(m.message)
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
