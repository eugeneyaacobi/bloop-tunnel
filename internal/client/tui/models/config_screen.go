package models

import (
	"fmt"
	"net/url"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// ConfigScreenSaveMsg is sent when user saves the config
type ConfigScreenSaveMsg struct {
	Config ConfigScreenData
}

// ConfigScreenCancelMsg is sent when user cancels
type ConfigScreenCancelMsg struct{}

// ConfigScreenSkipMsg is sent when user skips to next screen
type ConfigScreenSkipMsg struct {
	Config ConfigScreenData
}

// ConfigScreenData represents the configuration form data
type ConfigScreenData struct {
	ControlPlaneURL    string
	RelayURL           string
	AuthTokenEnv       string
	EnrollmentToken    string
	EnrollmentTokenEnv string
}

// ConfigScreenModel for global configuration
type ConfigScreenModel struct {
	data       ConfigScreenData
	fields     []InputFieldModel
	focused    int
	width      int
	height     int
	showAdvanced bool
}

// ConfigScreenOpts for constructing the model
type ConfigScreenOpts struct {
	Config        ConfigScreenData
	ShowAdvanced  bool
}

func NewConfigScreenModel(opts ConfigScreenOpts) ConfigScreenModel {
	data := opts.Config
	return ConfigScreenModel{
		data:       data,
		focused:    0,
		width:      80,
		height:     24,
		showAdvanced: opts.ShowAdvanced,
		fields: []InputFieldModel{
			NewInputField(InputFieldOpts{
				Label:       "Control Plane URL",
				Placeholder: "https://api.bloop.to",
				Value:       data.ControlPlaneURL,
				Validation:  validateURL,
			}),
			NewInputField(InputFieldOpts{
				Label:       "Relay URL",
				Placeholder: "wss://relay.bloop.to/connect",
				Value:       data.RelayURL,
				Validation:  validateWebSocketURL,
			}),
			NewInputField(InputFieldOpts{
				Label:       "Auth Token Env",
				Placeholder: "BLOOP_CLIENT_TOKEN",
				Value:       data.AuthTokenEnv,
			}),
		},
	}
}

func validateURL(v string) error {
	if v == "" {
		return fmt.Errorf("URL is required")
	}
	if _, err := url.Parse(v); err != nil {
		return fmt.Errorf("invalid URL format")
	}
	return nil
}

func validateWebSocketURL(v string) error {
	if v == "" {
		return fmt.Errorf("relay URL is required")
	}
	u, err := url.Parse(v)
	if err != nil {
		return fmt.Errorf("invalid URL format")
	}
	if u.Scheme != "ws" && u.Scheme != "wss" {
		return fmt.Errorf("must use ws:// or wss:// scheme")
	}
	return nil
}

func (m ConfigScreenModel) Init() tea.Cmd {
	return nil
}

func (m ConfigScreenModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "esc":
			return m, func() tea.Msg { return ConfigScreenCancelMsg{} }
		case "enter":
			if m.validate() {
				return m, func() tea.Msg { return ConfigScreenSaveMsg{Config: m.collectData()} }
			}
		case "a":
			m.showAdvanced = !m.showAdvanced
			return m, nil
		case "s":
			if m.validate() {
				return m, func() tea.Msg { return ConfigScreenSkipMsg{Config: m.collectData()} }
			}
		case "up", "shift+tab":
			m.blurFocused()
			m.focused = (m.focused - 1 + len(m.fields)) % len(m.fields)
			m.focusField()
		case "down", "tab":
			m.blurFocused()
			m.focused = (m.focused + 1) % len(m.fields)
			m.focusField()
		}

		// Forward to focused field
		if m.focused < len(m.fields) {
			m.fields[m.focused], _ = m.fields[m.focused].Update(msg)
		}

	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
	}

	return m, nil
}

func (m *ConfigScreenModel) focusField() {
	if m.focused < len(m.fields) {
		m.fields[m.focused] = m.fields[m.focused].Focused()
	}
}

func (m *ConfigScreenModel) blurFocused() {
	if m.focused < len(m.fields) {
		m.fields[m.focused] = m.fields[m.focused].Blur()
	}
}

func (m ConfigScreenModel) validate() bool {
	if validateURL(m.fields[0].GetValue()) != nil {
		return false
	}
	if validateWebSocketURL(m.fields[1].GetValue()) != nil {
		return false
	}
	return true
}

func (m ConfigScreenModel) collectData() ConfigScreenData {
	return ConfigScreenData{
		ControlPlaneURL: m.fields[0].GetValue(),
		RelayURL:        m.fields[1].GetValue(),
		AuthTokenEnv:    m.fields[2].GetValue(),
	}
}

func (m ConfigScreenModel) View() string {
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

	advancedStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FF9500")).
		Italic(true).
		MarginBottom(1)

	var fields []string
	for _, field := range m.fields {
		fields = append(fields, field.View())
	}

	advancedHint := ""
	if m.showAdvanced {
		advancedHint = advancedStyle.Render("[Advanced mode enabled - press A to disable]")
	} else {
		advancedHint = advancedStyle.Render("[Press A for advanced options]")
	}

	return lipgloss.JoinVertical(
		lipgloss.Left,
		titleStyle.Render("Global Configuration"),
		"",
		headerStyle.Render("Configure connection settings:"),
		"",
		advancedHint,
		"",
		lipgloss.JoinVertical(lipgloss.Left, fields...),
		"",
		footerStyle.Render("Tab: Navigate • Enter: Save & Continue • S: Skip with defaults • Q: Quit"),
	)
}
