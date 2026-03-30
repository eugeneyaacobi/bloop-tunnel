package models

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// OutputScreenConfirmMsg is sent when user confirms output settings
type OutputScreenConfirmMsg struct {
	Mode   string
	Path   string
}

// OutputScreenBackMsg is sent when user wants to go back
type OutputScreenBackMsg struct{}

// OutputScreenModel for selecting output format and location
type OutputScreenModel struct {
	formatField  SelectFieldModel
	pathField    InputFieldModel
	focused      int
	width        int
	height       int
	defaultPath  string
}

// OutputScreenOpts for constructing the model
type OutputScreenOpts struct {
	OutputMode  string
	OutputPath  string
}

func NewOutputScreenModel(opts OutputScreenOpts) OutputScreenModel {
	outputPath := opts.OutputPath
	if outputPath == "" {
		outputPath = "~/.bloop-tunnel/config.yaml"
	}

	modeIdx := 0
	switch opts.OutputMode {
	case "env-file":
		modeIdx = 1
	case "compose-block":
		modeIdx = 2
	}

	return OutputScreenModel{
		focused:     0,
		width:       80,
		height:      24,
		defaultPath: "~/.bloop-tunnel/config.yaml",
		formatField: NewSelectField(SelectFieldOpts{
			Label:    "Output Format",
			Options:  []string{"YAML File", "Environment File (.env)", "Docker Compose Block"},
			Selected: modeIdx,
		}),
		pathField: NewInputField(InputFieldOpts{
			Label:       "Output Path",
			Placeholder: "~/.bloop-tunnel/config.yaml",
			Value:       outputPath,
		}),
	}
}

func (m OutputScreenModel) Init() tea.Cmd {
	return nil
}

func (m OutputScreenModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "esc", "b":
			return m, func() tea.Msg { return OutputScreenBackMsg{} }
		case "enter":
			mode := m.selectedMode()
			path := m.pathField.GetValue()
			if path == "" {
				path = m.defaultPath
			}
			return m, func() tea.Msg { return OutputScreenConfirmMsg{Mode: mode, Path: path} }
		case "up", "shift+tab":
			m.blurFocused()
			m.focused = (m.focused - 1 + 2) % 2
			m.focusField()
		case "down", "tab":
			m.blurFocused()
			m.focused = (m.focused + 1) % 2
			m.focusField()
		}

		// Forward to focused field
		if m.focused == 0 {
			var cmd tea.Cmd
			m.formatField, cmd = m.formatField.Update(msg)
			return m, cmd
		} else if m.focused == 1 {
			var cmd tea.Cmd
			m.pathField, cmd = m.pathField.Update(msg)
			return m, cmd
		}

	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
	}

	return m, nil
}

func (m *OutputScreenModel) focusField() {
	if m.focused == 0 {
		m.formatField = m.formatField.Focused()
	} else if m.focused == 1 {
		m.pathField = m.pathField.Focused()
	}
}

func (m *OutputScreenModel) blurFocused() {
	if m.focused == 0 {
		m.formatField = m.formatField.Blur()
	} else if m.focused == 1 {
		m.pathField = m.pathField.Blur()
	}
}

func (m OutputScreenModel) selectedMode() string {
	switch m.formatField.SelectedIndex() {
	case 0:
		return "yaml"
	case 1:
		return "env-file"
	case 2:
		return "compose-block"
	default:
		return "yaml"
	}
}

func (m OutputScreenModel) View() string {
	titleStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#007AFF")).
		Bold(true).
		MarginTop(1).
		MarginBottom(1)

	headerStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#8E8E93")).
		Italic(true).
		MarginBottom(1)

	footerStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#8E8E93")).
		MarginTop(1)

	infoStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#8E8E93")).
		Italic(true).
		MarginLeft(4).
		MarginBottom(1)

	// Format descriptions
	formatDescriptions := map[string]string{
		"YAML File":           "Generate a YAML config file",
		"Environment File":    "Generate a .env file for Docker",
		"Docker Compose Block": "Generate a docker-compose.yml snippet",
	}

	selectedFormat := m.formatField.Selected()
	description := formatDescriptions[selectedFormat]

	return lipgloss.JoinVertical(
		lipgloss.Left,
		titleStyle.Render("Output Configuration"),
		"",
		headerStyle.Render("Choose how to generate your configuration:"),
		"",
		m.formatField.View(),
		infoStyle.Render("  → "+description),
		"",
		m.pathField.View(),
		"",
		footerStyle.Render("Tab: Navigate • Space/↑/↓: Change format • Enter: Generate • B: Back • Q: Quit"),
	)
}
