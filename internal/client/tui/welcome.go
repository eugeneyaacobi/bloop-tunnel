package tui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type WelcomeModel struct {
	choices       []WelcomeChoice
	selected      int
	width, height int
}

type WelcomeChoice struct {
	ID          string
	Label       string
	Description string
	Action      func() tea.Cmd
}

type WelcomeChoiceMsg struct{ Choice string }

func NewWelcomeModel() WelcomeModel {
	return WelcomeModel{
		choices: []WelcomeChoice{
			{
				ID:          "new",
				Label:       "New Setup",
				Description: "Configure bloop-tunnel from scratch with production defaults",
				Action:      func() tea.Cmd { return func() tea.Msg { return ScreenTransitionMsg{From: 0, To: 1} } },
			},
			{
				ID:          "edit",
				Label:       "Edit Existing",
				Description: "Load and edit an existing configuration file",
				Action:      func() tea.Cmd { return func() tea.Msg { return ScreenTransitionMsg{From: 0, To: 1} } },
			},
			{
				ID:          "quick",
				Label:       "Quick Start",
				Description: "Skip to output with default configuration",
				Action:      func() tea.Cmd { return func() tea.Msg { return ScreenTransitionMsg{From: 0, To: 6} } },
			},
		},
		selected: 0,
		width:    80,
		height:   24,
	}
}

func (m WelcomeModel) Init() tea.Cmd {
	return nil
}

func (m WelcomeModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			if m.selected > 0 {
				m.selected--
			}
		case "down", "j":
			if m.selected < len(m.choices)-1 {
				m.selected++
			}
		case "enter", " ":
			return m, m.choices[m.selected].Action()
		case "q", "ctrl+c":
			return m, tea.Quit
		}
	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
	}

	return m, nil
}

func (m WelcomeModel) View() string {
	menuWidth := min(50, m.width-8)

	choiceStyle := lipgloss.NewStyle().
		Padding(0, 2).
		MarginBottom(1).
		Width(menuWidth)

	selectedStyle := lipgloss.NewStyle().
		Padding(0, 2).
		MarginBottom(1).
		Width(menuWidth).
		Background(primaryColor).
		Foreground(textColor).
		Bold(true)

	var choices []string
	for i, choice := range m.choices {
		if i == m.selected {
			choices = append(choices, selectedStyle.Render(fmt.Sprintf("• %s", choice.Label)))
		} else {
			choices = append(choices, choiceStyle.Render(fmt.Sprintf("  %s", choice.Label)))
		}
		choices = append(choices, descStyle.Render("    "+choice.Description))
	}

	return lipgloss.JoinVertical(
		lipgloss.Left,
		WelcomeBanner,
		titleStyle.Render("Welcome to bloop-tunnel Setup"),
		subtitleStyle.Render("Choose how you'd like to configure your tunnels"),
		"",
		lipgloss.JoinVertical(lipgloss.Left, choices...),
		"",
		footerStyle.Render("↑/↓: Navigate • Enter: Select • q: Quit"),
	)
}
