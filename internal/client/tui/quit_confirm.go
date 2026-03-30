package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type QuitConfirmModel struct {
	message       string
	options       []string
	selected      int
	width, height int
}

func NewQuitConfirmModel() QuitConfirmModel {
	return QuitConfirmModel{
		message:  "You have unsaved changes. Are you sure you want to quit?",
		options:  []string{"Quit", "Cancel"},
		selected: 1, // Default to Cancel
		width:    80,
		height:   24,
	}
}

func (m QuitConfirmModel) Init() tea.Cmd {
	return nil
}

func (m QuitConfirmModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "left", "h":
			if m.selected > 0 {
				m.selected--
			}
		case "right", "l":
			if m.selected < len(m.options)-1 {
				m.selected++
			}
		case "tab":
			m.selected = (m.selected + 1) % len(m.options)
		case "enter", " ":
			return m, func() tea.Msg { return QuitConfirmMsg{Choice: m.options[m.selected]} }
		case "esc", "q":
			return m, func() tea.Msg { return QuitConfirmMsg{Choice: "Cancel"} }
		case "ctrl+c":
			return m, tea.Quit
		}
	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
	}

	return m, nil
}

func (m QuitConfirmModel) View() string {
	messageSty := lipgloss.NewStyle().Foreground(textColor).MarginBottom(2)

	var options []string
	for i, opt := range m.options {
		if i == m.selected {
			options = append(options, optionSelectedStyle.Render(opt))
		} else {
			options = append(options, optionStyle.Render(opt))
		}
	}

	boxContent := lipgloss.JoinVertical(
		lipgloss.Center,
		messageSty.Render(m.message),
		"",
		lipgloss.JoinHorizontal(lipgloss.Center, options...),
	)

	box := dialogBoxStyle.Render(boxContent)

	// Guard against zero dimensions
	if m.width > 0 && m.height > 0 {
		return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, box)
	}
	return box
}
