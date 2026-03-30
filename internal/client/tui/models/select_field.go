package models

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type SelectFieldModel struct {
	label       string
	options     []string
	selected    int
	focused     bool
}

type SelectFieldOpts struct {
	Label       string
	Options     []string
	Selected    int
}

type SelectFieldChangeMsg struct{ Index int }

func NewSelectField(opts SelectFieldOpts) SelectFieldModel {
	selected := opts.Selected
	if selected < 0 || selected >= len(opts.Options) {
		selected = 0
	}

	return SelectFieldModel{
		label:    opts.Label,
		options:  opts.Options,
		selected: selected,
		focused:  false,
	}
}

func (m SelectFieldModel) Init() tea.Cmd {
	return nil
}

func (m SelectFieldModel) Update(msg tea.Msg) (SelectFieldModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			if m.selected > 0 {
				m.selected--
			}
		case "down", "j":
			if m.selected < len(m.options)-1 {
				m.selected++
			}
		case "home", "g":
			m.selected = 0
		case "end", "G":
			m.selected = len(m.options) - 1
		case "space":
			// Cycle to next option
			m.selected = (m.selected + 1) % len(m.options)
			return m, func() tea.Msg { return SelectFieldChangeMsg{Index: m.selected} }
		case "enter", " ":
			return m, func() tea.Msg { return SelectFieldChangeMsg{Index: m.selected} }
		}
	case SelectFieldChangeMsg:
		m.selected = msg.Index
	}

	return m, nil
}

func (m SelectFieldModel) View() string {
	labelStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#007AFF")).
		Bold(true).
		MarginBottom(0)

	optionStyle := lipgloss.NewStyle().
		Padding(0, 1).
		MarginRight(1)

	selectedStyle := optionStyle.Copy().
		Foreground(lipgloss.Color("#007AFF")).
		Bold(true)

	var rows []string
	for i, opt := range m.options {
		prefix := "  "
		if i == m.selected && m.focused {
			prefix = "● "
			selectedStyle = selectedStyle.Copy().
				Background(lipgloss.Color("#007AFF")).
				Foreground(lipgloss.Color("#FFFFFF"))
		} else if i == m.selected {
			prefix = "○ "
		}
		rows = append(rows, prefix+selectedStyle.Render(opt))
	}

	boxStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#007AFF")).
		Padding(0, 1).
		Width(52)

	if m.focused {
		boxStyle = boxStyle.Copy().
			Foreground(lipgloss.Color("#007AFF")).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#007AFF")).
			Bold(true)
	}

	label := labelStyle.Render(m.label + ":")

	return lipgloss.JoinVertical(
		lipgloss.Left,
		boxStyle.Render(lipgloss.JoinVertical(lipgloss.Left, label, lipgloss.JoinVertical(lipgloss.Left, rows...))),
	)
}

func (m SelectFieldModel) WithSelection(index int) SelectFieldModel {
	m.selected = index
	return m
}

func (m SelectFieldModel) WithOptions(options []string) SelectFieldModel {
	m.options = options
	if m.selected >= len(m.options) {
		m.selected = 0
	}
	return m
}

func (m SelectFieldModel) SelectedIndex() int {
	return m.selected
}

func (m SelectFieldModel) Selected() string {
	if m.selected >= 0 && m.selected < len(m.options) {
		return m.options[m.selected]
	}
	return ""
}

func (m SelectFieldModel) Focused() SelectFieldModel {
	m.focused = true
	return m
}

func (m SelectFieldModel) Blur() SelectFieldModel {
	m.focused = false
	return m
}
