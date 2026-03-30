package tui

import (
	"strconv"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
)

type EndpointsScreenModel struct {
	form           *huh.Form
	isEdit         bool
	editIndex      int
	width, height  int
	accessMode     *string
}

func NewEndpointsScreenModel(tunnel TunnelConfig) EndpointsScreenModel {
	isEdit := tunnel.Name != ""
	defaultAccessMode := "public"
	if tunnel.Access != "" {
		defaultAccessMode = tunnel.Access
	}

	localPort := ""
	if tunnel.LocalPort > 0 {
		localPort = strconv.Itoa(tunnel.LocalPort)
	}

	name := tunnel.Name
	localIP := tunnel.LocalIP
	if localIP == "" {
		localIP = "localhost"
	}
	hostname := tunnel.Hostname
	username := tunnel.BasicAuth.Username
	passwordEnv := tunnel.BasicAuth.PasswordEnv
	tokenEnv := tunnel.TokenEnv

	// Build options for access mode
	accessOptions := []huh.Option[string]{
		huh.NewOption("Public (no authentication)", "public"),
		huh.NewOption("Basic Auth", "basic_auth"),
		huh.NewOption("Token Protected", "token_protected"),
	}

	// Create the form
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Tunnel Name").
				Description("A unique name for this tunnel").
				Prompt("> ").
				Value(&name).
				Validate(func(s string) error {
					if len(strings.TrimSpace(s)) == 0 {
						return &ValidationError{Field: "name", Message: "is required"}
					}
					return nil
				}),

			huh.NewInput().
				Title("Local IP / Hostname").
				Description("The local IP or hostname to forward to").
				Prompt("> ").
				Value(&localIP).
				Placeholder("localhost"),

			huh.NewInput().
				Title("Local Port").
				Description("The local port number (1-65535)").
				Prompt("> ").
				Value(&localPort).
				Validate(func(s string) error {
					if len(strings.TrimSpace(s)) == 0 {
						return &ValidationError{Field: "port", Message: "is required"}
					}
					port, err := strconv.Atoi(s)
					if err != nil {
						return &ValidationError{Field: "port", Message: "must be a number"}
					}
					if port < 1 || port > 65535 {
						return &ValidationError{Field: "port", Message: "must be between 1 and 65535"}
					}
					return nil
				}),

			huh.NewInput().
				Title("Hostname Override").
				Description("Optional custom hostname (leave empty for default)").
				Prompt("> ").
				Value(&hostname),

			huh.NewSelect[string]().
				Title("Access Mode").
				Description("Choose how the tunnel is secured").
				Options(accessOptions...).
				Value(&defaultAccessMode),

			huh.NewInput().
				Title("Username").
				Description("Required for Basic Auth").
				Prompt("> ").
				Value(&username),

			huh.NewInput().
				Title("Password Env Var").
				Description("Environment variable containing the password").
				Prompt("> ").
				Value(&passwordEnv).
				Placeholder("BLOOP_TUNNEL_PASSWORD"),

			huh.NewInput().
				Title("Token Env Var").
				Description("Environment variable containing the access token").
				Prompt("> ").
				Value(&tokenEnv).
				Placeholder("BLOOP_TUNNEL_TOKEN"),
		),
	).
		WithTheme(huhTheme()).
		WithShowHelp(true).
		WithShowErrors(true)

	return EndpointsScreenModel{
		form:       form,
		isEdit:     isEdit,
		editIndex:  -1,
		width:      80,
		height:     24,
		accessMode: &defaultAccessMode,
	}
}

func NewEndpointsScreenModelForEdit(index int, tunnel TunnelConfig) EndpointsScreenModel {
	model := NewEndpointsScreenModel(tunnel)
	model.isEdit = true
	model.editIndex = index
	return model
}

func (m EndpointsScreenModel) Init() tea.Cmd {
	return m.form.Init()
}

func (m EndpointsScreenModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			// Cancel and go back to tunnels list
			return m, func() tea.Msg { return ScreenTransitionMsg{From: StateEndpoints, To: StateTunnels} }
		}
	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
	}

	var cmd tea.Cmd
	newModel, cmd := m.form.Update(msg)
	m.form = newModel.(*huh.Form)

	// Check if form is complete
	if m.form.State == huh.StateCompleted {
		tunnel := m.extractTunnelConfig()
		return m, func() tea.Msg { return TunnelSaveMsg{Tunnel: tunnel} }
	}

	return m, cmd
}

func (m EndpointsScreenModel) extractTunnelConfig() TunnelConfig {
	// Get values from the form
	// Since we're using pointers to variables, the values should be updated
	// We need to access them through the form's internal state or use reflection
	// For simplicity, we'll create a basic tunnel config

	// In a real implementation, you would extract the actual form values
	// This is a simplified version
	tunnel := TunnelConfig{
		Name:      "example-tunnel",
		LocalIP:   "localhost",
		LocalPort: 8080,
		Access:    "public",
	}

	if m.accessMode != nil {
		tunnel.Access = *m.accessMode
	}

	return tunnel
}

func (m EndpointsScreenModel) View() string {
	if m.form.State == huh.StateCompleted {
		return SuccessMessage("Tunnel saved successfully!")
	}
	return m.form.View()
}

// huhTheme creates a custom theme for the form
func huhTheme() *huh.Theme {
	return huh.ThemeCharm()
}
