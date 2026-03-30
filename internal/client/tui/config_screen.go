package tui

import (
	"strconv"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"bloop-tunnel/internal/client/tui/models"
)

type ConfigScreenModel struct {
	fields            []models.InputFieldModel
	focusedField      int
	config            Config
	width, height     int
	showHelp          bool
}

func NewConfigScreenModel(config Config) ConfigScreenModel {
	initialDelay := strconv.Itoa(config.Reconnect.InitialDelayMs)
	maxDelay := strconv.Itoa(config.Reconnect.MaxDelayMs)

	fields := []models.InputFieldModel{
		models.NewInputField(models.InputFieldOpts{
			Label:       "Control Plane URL",
			Placeholder: "https://api.bloop.to",
			Value:       config.ControlPlaneURL,
		}),
		models.NewInputField(models.InputFieldOpts{
			Label:       "Relay URL",
			Placeholder: "wss://relay.bloop.to/connect",
			Value:       config.RelayURL,
		}),
		models.NewInputField(models.InputFieldOpts{
			Label:       "Tunnel Token",
			Placeholder: "Paste your tunnel token from bloop.to",
			Value:       config.AuthTokenEnv,
		}),
		models.NewInputField(models.InputFieldOpts{
			Label:       "Enrollment Token",
			Placeholder: "(optional)",
			Value:       config.EnrollmentToken,
		}),
		models.NewInputField(models.InputFieldOpts{
			Label:       "Enrollment Token Env Var",
			Placeholder: "(optional)",
			Value:       config.EnrollmentTokenEnv,
		}),
		models.NewInputField(models.InputFieldOpts{
			Label:       "Reconnect Initial Delay (ms)",
			Placeholder: "1000",
			Value:       initialDelay,
			Validation: func(s string) error {
				if s == "" {
					return nil
				}
				val, err := strconv.Atoi(s)
				if err != nil {
					return err
				}
				if val < 0 || val > 60000 {
					return &ValidationError{Field: "initial delay", Message: "must be between 0 and 60000"}
				}
				return nil
			},
		}),
		models.NewInputField(models.InputFieldOpts{
			Label:       "Reconnect Max Delay (ms)",
			Placeholder: "30000",
			Value:       maxDelay,
			Validation: func(s string) error {
				if s == "" {
					return nil
				}
				val, err := strconv.Atoi(s)
				if err != nil {
					return err
				}
				if val < 0 || val > 300000 {
					return &ValidationError{Field: "max delay", Message: "must be between 0 and 300000"}
				}
				return nil
			},
		}),
	}

	// Focus first field
	fields[0] = fields[0].Focused()

	return ConfigScreenModel{
		fields:       fields,
		focusedField: 0,
		config:       config,
		width:        80,
		height:       24,
		showHelp:     false,
	}
}

type ValidationError struct {
	Field   string
	Message string
}

func (e *ValidationError) Error() string {
	return e.Field + ": " + e.Message
}

func (m ConfigScreenModel) Init() tea.Cmd {
	return nil
}

func (m ConfigScreenModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "tab":
			// Move to next field
			m.blurFocusedField()
			m.focusedField = (m.focusedField + 1) % len(m.fields)
			m.focusField()
			return m, nil
		case "shift+tab":
			// Move to previous field
			m.blurFocusedField()
			m.focusedField = (m.focusedField - 1 + len(m.fields)) % len(m.fields)
			m.focusField()
			return m, nil
		case "enter":
			// Save and advance
			return m, m.saveConfig()
		case "esc", "ctrl+c":
			// Go back to welcome
			return m, func() tea.Msg { return ScreenTransitionMsg{From: StateConfig, To: StateWelcome} }
		case "?":
			m.showHelp = !m.showHelp
			return m, nil
		}
	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
	}

	// Delegate update to focused field
	var cmd tea.Cmd
	m.fields[m.focusedField], cmd = m.fields[m.focusedField].Update(msg)
	return m, cmd
}

func (m *ConfigScreenModel) focusField() {
	m.fields[m.focusedField] = m.fields[m.focusedField].Focused()
}

func (m *ConfigScreenModel) blurFocusedField() {
	m.fields[m.focusedField] = m.fields[m.focusedField].Blur()
}

func (m ConfigScreenModel) saveConfig() tea.Cmd {
	// Collect values from fields
	m.config.ControlPlaneURL = m.fields[0].View()
	m.config.RelayURL = m.fields[1].View()
	m.config.AuthTokenEnv = m.fields[2].View()
	m.config.EnrollmentToken = m.fields[3].View()
	m.config.EnrollmentTokenEnv = m.fields[4].View()

	// Parse numeric fields
	initialDelay, _ := strconv.Atoi(m.fields[5].View())
	maxDelay, _ := strconv.Atoi(m.fields[6].View())
	m.config.Reconnect.InitialDelayMs = initialDelay
	m.config.Reconnect.MaxDelayMs = maxDelay

	return func() tea.Msg {
		return ConfigSaveMsg{Config: m.config}
	}
}

func (m ConfigScreenModel) View() string {
	title := titleStyle.Render("Configuration Settings")
	subtitle := subtitleStyle.Render("Configure your bloop-tunnel connection settings")

	var fieldViews []string
	for i, field := range m.fields {
		fieldViews = append(fieldViews, field.View())
		if i < len(m.fields)-1 {
			fieldViews = append(fieldViews, "")
		}
	}

	fieldsContent := lipgloss.JoinVertical(lipgloss.Left, fieldViews...)

	content := lipgloss.JoinVertical(
		lipgloss.Left,
		title,
		subtitle,
		"",
		fieldsContent,
	)

	// Help footer
	var helpText string
	if m.showHelp {
		helpText = HelpText("Tab: Next field • Shift+Tab: Prev field • Enter: Save • Esc: Back")
	} else {
		helpText = footerStyle.Render("?: Show help")
	}

	content = lipgloss.JoinVertical(lipgloss.Left, content, "", helpText)

	return Center(m.width, content)
}
