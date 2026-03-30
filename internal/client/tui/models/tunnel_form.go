package models

import (
	"fmt"
	"strconv"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// TunnelFormSaveMsg is sent when user saves the tunnel form
type TunnelFormSaveMsg struct {
	Tunnel TunnelFormData
}

// TunnelFormCancelMsg is sent when user cancels the form
type TunnelFormCancelMsg struct{}

// TunnelFormData represents the form data
type TunnelFormData struct {
	Name       string
	Hostname   string
	LocalIP    string
	LocalPort  string
	Access     string
	BasicAuth  BasicAuthFormData
	TokenEnv   string
}

type BasicAuthFormData struct {
	Username    string
	PasswordEnv string
}

// TunnelFormModel for editing/creating tunnel configuration
type TunnelFormModel struct {
	data        TunnelFormData
	fields      []InputFieldModel
	accessField SelectFieldModel
	basicAuthFields []InputFieldModel
	tokenField  InputFieldModel
	focused     int
	width       int
	height      int
	isEditing   bool
}

// TunnelFormOpts for constructing the model
type TunnelFormOpts struct {
	Tunnel     TunnelFormData
	IsEditing  bool
}

func NewTunnelFormModel(opts TunnelFormOpts) TunnelFormModel {
	data := opts.Tunnel
	if data.Access == "" {
		data.Access = "public"
	}

	return TunnelFormModel{
		data:      data,
		focused:   0,
		width:     80,
		height:    24,
		isEditing: opts.IsEditing,
		fields: []InputFieldModel{
			NewInputField(InputFieldOpts{
				Label:       "Name",
				Placeholder: "my-tunnel",
				Value:       data.Name,
				Validation:  validateTunnelName,
			}),
			NewInputField(InputFieldOpts{
				Label:       "Hostname",
				Placeholder: "optional.bloop.to",
				Value:       data.Hostname,
			}),
			NewInputField(InputFieldOpts{
				Label:       "Local IP",
				Placeholder: "localhost",
				Value:       data.LocalIP,
			}),
			NewInputField(InputFieldOpts{
				Label:       "Local Port",
				Placeholder: "8080",
				Value:       data.LocalPort,
				Validation:  validatePort,
			}),
		},
		accessField: NewSelectField(SelectFieldOpts{
			Label:    "Access Control",
			Options:  []string{"public", "basic_auth", "token_protected"},
			Selected: accessToIndex(data.Access),
		}),
		basicAuthFields: []InputFieldModel{
			NewInputField(InputFieldOpts{
				Label:       "Username",
				Placeholder: "admin",
				Value:       data.BasicAuth.Username,
			}),
			NewInputField(InputFieldOpts{
				Label:       "Password Env",
				Placeholder: "TUNNEL_PASSWORD",
				Value:       data.BasicAuth.PasswordEnv,
			}),
		},
		tokenField: NewInputField(InputFieldOpts{
			Label:       "Token Env",
			Placeholder: "TUNNEL_TOKEN",
			Value:       data.TokenEnv,
		}),
	}
}

func accessToIndex(access string) int {
	switch access {
	case "public":
		return 0
	case "basic_auth":
		return 1
	case "token_protected":
		return 2
	default:
		return 0
	}
}

func indexToAccess(index int) string {
	switch index {
	case 0:
		return "public"
	case 1:
		return "basic_auth"
	case 2:
		return "token_protected"
	default:
		return "public"
	}
}

func validateTunnelName(v string) error {
	if v == "" {
		return fmt.Errorf("name is required")
	}
	if len(v) > 63 {
		return fmt.Errorf("name too long (max 63 characters)")
	}
	return nil
}

func validatePort(v string) error {
	if v == "" {
		return fmt.Errorf("port is required")
	}
	port, err := strconv.Atoi(v)
	if err != nil {
		return fmt.Errorf("invalid port number")
	}
	if port < 1 || port > 65535 {
		return fmt.Errorf("port must be between 1 and 65535")
	}
	return nil
}

func (m TunnelFormModel) Init() tea.Cmd {
	return nil
}

func (m TunnelFormModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "esc", "b":
			return m, func() tea.Msg { return TunnelFormCancelMsg{} }
		case "enter":
			if m.focused == m.totalFields()-1 {
				// Save and submit
				if m.validate() {
					return m, func() tea.Msg { return TunnelFormSaveMsg{Tunnel: m.collectData()} }
				}
			} else {
				// Move to next field
				m.blurFocused()
				m.focused = (m.focused + 1) % m.totalFields()
				m.focusField()
			}
		case "up", "shift+tab":
			m.blurFocused()
			m.focused = (m.focused - 1 + m.totalFields()) % m.totalFields()
			m.focusField()
		case "down", "tab":
			m.blurFocused()
			m.focused = (m.focused + 1) % m.totalFields()
			m.focusField()
		}

		// Forward to focused field
		m.updateFocusedField(msg)

	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
	}

	return m, nil
}

func (m *TunnelFormModel) totalFields() int {
	baseFields := len(m.fields) + 1 // +1 for access field
	accessType := indexToAccess(m.accessField.SelectedIndex())
	if accessType == "basic_auth" {
		return baseFields + len(m.basicAuthFields)
	} else if accessType == "token_protected" {
		return baseFields + 1
	}
	return baseFields
}

func (m *TunnelFormModel) focusField() {
	baseFields := len(m.fields)
	accessType := indexToAccess(m.accessField.SelectedIndex())

	if m.focused < baseFields {
		m.fields[m.focused] = m.fields[m.focused].Focused()
	} else if m.focused == baseFields {
		m.accessField = m.accessField.Focused()
	} else if accessType == "basic_auth" {
		idx := m.focused - baseFields - 1
		if idx < len(m.basicAuthFields) {
			m.basicAuthFields[idx] = m.basicAuthFields[idx].Focused()
		}
	} else if accessType == "token_protected" {
		m.tokenField = m.tokenField.Focused()
	}
}

func (m *TunnelFormModel) blurFocused() {
	baseFields := len(m.fields)
	accessType := indexToAccess(m.accessField.SelectedIndex())

	if m.focused < baseFields {
		m.fields[m.focused] = m.fields[m.focused].Blur()
	} else if m.focused == baseFields {
		m.accessField = m.accessField.Blur()
	} else if accessType == "basic_auth" {
		idx := m.focused - baseFields - 1
		if idx < len(m.basicAuthFields) {
			m.basicAuthFields[idx] = m.basicAuthFields[idx].Blur()
		}
	} else if accessType == "token_protected" {
		m.tokenField = m.tokenField.Blur()
	}
}

func (m *TunnelFormModel) updateFocusedField(msg tea.Msg) {
	baseFields := len(m.fields)
	accessType := indexToAccess(m.accessField.SelectedIndex())

	if m.focused < baseFields {
		
		m.fields[m.focused], _ = m.fields[m.focused].Update(msg)
	} else if m.focused == baseFields {
		
		m.accessField, _ = m.accessField.Update(msg)
		{
			// Access type changed, reset focus
			m.blurFocused()
			m.focused = baseFields
			m.focusField()
		}
	} else if accessType == "basic_auth" {
		idx := m.focused - baseFields - 1
		if idx < len(m.basicAuthFields) {
			
			m.basicAuthFields[idx], _ = m.basicAuthFields[idx].Update(msg)
		}
	} else if accessType == "token_protected" {
		
		m.tokenField, _ = m.tokenField.Update(msg)
	}
}

func (m TunnelFormModel) validate() bool {
	// Validate name
	if validateTunnelName(m.fields[0].GetValue()) != nil {
		return false
	}
	// Validate port
	if validatePort(m.fields[3].GetValue()) != nil {
		return false
	}
	return true
}

func (m TunnelFormModel) collectData() TunnelFormData {
	data := TunnelFormData{
		Name:      m.fields[0].GetValue(),
		Hostname:  m.fields[1].GetValue(),
		LocalIP:   m.fields[2].GetValue(),
		LocalPort: m.fields[3].GetValue(),
		Access:    indexToAccess(m.accessField.SelectedIndex()),
	}

	if data.Access == "basic_auth" {
		data.BasicAuth = BasicAuthFormData{
			Username:    m.basicAuthFields[0].GetValue(),
			PasswordEnv: m.basicAuthFields[1].GetValue(),
		}
	} else if data.Access == "token_protected" {
		data.TokenEnv = m.tokenField.GetValue()
	}

	return data
}

func (m TunnelFormModel) View() string {
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

	title := "Add Tunnel"
	if m.isEditing {
		title = "Edit Tunnel"
	}

	var fields []string

	// Base fields
	for _, field := range m.fields {
		fields = append(fields, field.View())
	}

	// Access field
	fields = append(fields, m.accessField.View())

	// Conditional fields based on access type
	accessType := indexToAccess(m.accessField.SelectedIndex())
	if accessType == "basic_auth" {
		for _, field := range m.basicAuthFields {
			fields = append(fields, field.View())
		}
	} else if accessType == "token_protected" {
		fields = append(fields, m.tokenField.View())
	}

	return lipgloss.JoinVertical(
		lipgloss.Left,
		titleStyle.Render(title),
		"",
		headerStyle.Render("Configure your tunnel endpoint:"),
		"",
		lipgloss.JoinVertical(lipgloss.Left, fields...),
		"",
		footerStyle.Render("Tab: Navigate • Enter: Save/Next • Esc: Cancel • Q: Quit"),
	)
}
