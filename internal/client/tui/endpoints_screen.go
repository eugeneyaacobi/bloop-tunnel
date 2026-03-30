package tui

import (
	"strconv"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
)

type EndpointsScreenModel struct {
	form          *huh.Form
	isEdit        bool
	editIndex     int
	width, height int

	// Value pointers — huh writes into these
	name        string
	localIP     string
	localPort   string
	hostname    string
	accessMode  string
	username    string
	passwordEnv string
	tokenEnv    string
}

func NewEndpointsScreenModel(tunnel TunnelConfig) EndpointsScreenModel {
	m := EndpointsScreenModel{
		isEdit:     tunnel.Name != "",
		name:       tunnel.Name,
		localIP:    tunnel.LocalIP,
		localPort:  strconv.Itoa(tunnel.LocalPort),
		hostname:   tunnel.Hostname,
		accessMode: tunnel.Access,
		username:   tunnel.BasicAuth.Username,
		passwordEnv: tunnel.BasicAuth.PasswordEnv,
		tokenEnv:   tunnel.TokenEnv,
		width:      80,
		height:     24,
		editIndex:  -1,
	}

	if m.localIP == "" {
		m.localIP = "localhost"
	}
	if m.accessMode == "" {
		m.accessMode = "public"
	}
	if tunnel.LocalPort > 0 {
		m.localPort = strconv.Itoa(tunnel.LocalPort)
	} else {
		m.localPort = ""
	}

	m.buildForm()
	return m
}

func NewEndpointsScreenModelForEdit(index int, tunnel TunnelConfig) EndpointsScreenModel {
	m := NewEndpointsScreenModel(tunnel)
	m.isEdit = true
	m.editIndex = index
	return m
}

func (m *EndpointsScreenModel) buildForm() {
	accessOptions := []huh.Option[string]{
		huh.NewOption("Public (no authentication)", "public"),
		huh.NewOption("Basic Auth", "basic_auth"),
		huh.NewOption("Token Protected", "token_protected"),
	}

	m.form = huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Tunnel Name").
				Description("A unique name for this tunnel").
				Prompt("> ").
				Value(&m.name).
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
				Value(&m.localIP).
				Placeholder("localhost"),

			huh.NewInput().
				Title("Local Port").
				Description("The local port number (1-65535)").
				Prompt("> ").
				Value(&m.localPort).
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
				Value(&m.hostname),

			huh.NewSelect[string]().
				Title("Access Mode").
				Description("Choose how the tunnel is secured").
				Options(accessOptions...).
				Value(&m.accessMode),

			huh.NewInput().
				Title("Username").
				Description("Required for Basic Auth").
				Prompt("> ").
				Value(&m.username),

			huh.NewInput().
				Title("Password Env Var").
				Description("Environment variable containing the password").
				Prompt("> ").
				Value(&m.passwordEnv).
				Placeholder("BLOOP_TUNNEL_PASSWORD"),

			huh.NewInput().
				Title("Token Env Var").
				Description("Environment variable containing the access token").
				Prompt("> ").
				Value(&m.tokenEnv).
				Placeholder("BLOOP_TUNNEL_TOKEN"),
		),
	).
		WithTheme(huhTheme()).
		WithShowHelp(true).
		WithShowErrors(true)
}

func (m EndpointsScreenModel) Init() tea.Cmd {
	return m.form.Init()
}

func (m EndpointsScreenModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			return m, func() tea.Msg { return ScreenTransitionMsg{From: StateEndpoints, To: StateTunnels} }
		}
	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
	}

	var cmd tea.Cmd
	newModel, cmd := m.form.Update(msg)
	m.form = newModel.(*huh.Form)

	if m.form.State == huh.StateCompleted {
		tunnel := m.extractTunnelConfig()
		return m, func() tea.Msg { return TunnelSaveMsg{Tunnel: tunnel} }
	}

	return m, cmd
}

func (m EndpointsScreenModel) extractTunnelConfig() TunnelConfig {
	port, _ := strconv.Atoi(m.localPort)
	tunnel := TunnelConfig{
		Name:      m.name,
		LocalIP:   m.localIP,
		LocalPort: port,
		Hostname:  m.hostname,
		Access:    m.accessMode,
		TokenEnv:  m.tokenEnv,
	}
	if m.accessMode == "basic_auth" {
		tunnel.BasicAuth = BasicAuthConfig{
			Username:    m.username,
			PasswordEnv: m.passwordEnv,
		}
	}
	return tunnel
}

func (m EndpointsScreenModel) View() string {
	if m.form.State == huh.StateCompleted {
		return SuccessMessage("Tunnel saved successfully!")
	}
	return m.form.View()
}

func huhTheme() *huh.Theme {
	return huh.ThemeCharm()
}
