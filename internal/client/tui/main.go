package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"bloop-tunnel/internal/client/tui/models"
)

type MainModel struct {
	state              State
	stateStack         []State
	config             Config
	currentScreen      tea.Model
	editingTunnelIndex int // Track which tunnel is being edited (-1 for new)
	width              int
	height             int
}

func NewMainModel() *MainModel {
	return &MainModel{
		state:              StateWelcome,
		stateStack:         []State{},
		config:             DefaultConfig(),
		editingTunnelIndex: -1,
		width:              80,
		height:             24,
	}
}

func (m *MainModel) Init() tea.Cmd {
	m.currentScreen = NewWelcomeModel()
	return nil
}

func (m *MainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			if m.config.HasChanges {
				m.state = StateQuitConfirm
				m.stateStack = append(m.stateStack, m.state)
				m.currentScreen = NewQuitConfirmModel()
				return m, nil
			}
			return m, tea.Quit
		}
	case ScreenTransitionMsg:
		m.state = msg.To
		m.currentScreen = m.getScreenForState(msg.To)
		if msg.From != msg.To {
			m.stateStack = append(m.stateStack, msg.From)
		}
		return m, nil
	case ConfigSaveMsg:
		m.config = msg.Config
		m.config.HasChanges = false
		// After config save, advance to tunnels
		if m.state == StateConfig {
			return m, func() tea.Msg { return ScreenTransitionMsg{From: StateConfig, To: StateTunnels} }
		}
		return m, nil
	case ConfigChangedMsg:
		m.config.HasChanges = true
		return m, nil
	case TunnelAddMsg:
		m.editingTunnelIndex = -1
		return m, func() tea.Msg { return ScreenTransitionMsg{From: StateTunnels, To: StateEndpoints} }
	case TunnelSaveMsg:
		m.saveTunnel(msg.Tunnel)
		m.config.HasChanges = true
		return m, func() tea.Msg { return ScreenTransitionMsg{From: StateEndpoints, To: StateTunnels} }
	case TunnelDeleteMsg:
		if msg.Index >= 0 && msg.Index < len(m.config.Tunnels) {
			m.config.Tunnels = append(m.config.Tunnels[:msg.Index], m.config.Tunnels[msg.Index+1:]...)
		}
		m.config.HasChanges = true
		// Refresh the tunnels screen
		m.currentScreen = NewTunnelsScreenModel(m.config.Tunnels)
		return m, nil
	case TunnelsListAddMsg:
		m.editingTunnelIndex = -1
		return m, func() tea.Msg { return ScreenTransitionMsg{From: StateTunnels, To: StateEndpoints} }
	case TunnelsListEditMsg:
		m.editingTunnelIndex = msg.Index
		return m, func() tea.Msg { return ScreenTransitionMsg{From: StateTunnels, To: StateEndpoints} }
	case TunnelsListNextMsg:
		return m, func() tea.Msg { return ScreenTransitionMsg{From: StateTunnels, To: StateReview} }
	case TunnelsListBackMsg:
		return m, func() tea.Msg { return ScreenTransitionMsg{From: StateTunnels, To: StateConfig} }
	case QuitConfirmMsg:
		if msg.Choice == "Quit" {
			return m, tea.Quit
		}
		if len(m.stateStack) > 0 {
			prevState := m.stateStack[len(m.stateStack)-1]
			m.stateStack = m.stateStack[:len(m.stateStack)-1]
			m.state = prevState
			m.currentScreen = m.getScreenForState(prevState)
			return m, nil
		}
	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
	}

	if m.currentScreen != nil {
		var cmd tea.Cmd
		m.currentScreen, cmd = m.currentScreen.Update(msg)
		return m, cmd
	}

	return m, nil
}

func (m *MainModel) saveTunnel(tunnel TunnelConfig) {
	if m.editingTunnelIndex >= 0 && m.editingTunnelIndex < len(m.config.Tunnels) {
		// Update existing tunnel
		m.config.Tunnels[m.editingTunnelIndex] = tunnel
	} else {
		// Add new tunnel
		m.config.Tunnels = append(m.config.Tunnels, tunnel)
	}
	m.editingTunnelIndex = -1
}

func (m MainModel) getScreenForState(state State) tea.Model {
	switch state {
	case StateWelcome:
		return NewWelcomeModel()
	case StateConfig:
		return NewConfigScreenModel(m.config)
	case StateEndpoints:
		if m.editingTunnelIndex >= 0 && m.editingTunnelIndex < len(m.config.Tunnels) {
			return NewEndpointsScreenModelForEdit(m.editingTunnelIndex, m.config.Tunnels[m.editingTunnelIndex])
		}
		return NewEndpointsScreenModel(TunnelConfig{})
	case StateTunnels:
		return NewTunnelsScreenModel(m.config.Tunnels)
	case StateVerification:
		return NewVerificationScreenModel(m.config)
	case StateDockerDiscovery:
		return NewDockerDiscoveryScreenModel(m.config)
	case StateReview:
		return NewReviewScreenModel(m.config)
	case StateOutput:
		return NewOutputScreenModel(m.config)
	case StateQuitConfirm:
		return NewQuitConfirmModel()
	default:
		return models.NewStatusModel()
	}
}

func (m *MainModel) View() string {
	if m.currentScreen != nil {
		return m.currentScreen.View()
	}
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FF3B30")).
		Render("Error: No screen active")
}

func DefaultConfig() Config {
	return Config{
		ControlPlaneURL: "https://api.bloop.to",
		RelayURL:        "wss://relay.bloop.to/connect",
		AuthTokenEnv:    "BLOOP_CLIENT_TOKEN",
		Reconnect: ReconnectConfig{
			InitialDelayMs: 1000,
			MaxDelayMs:     30000,
		},
		Tunnels:    []TunnelConfig{},
		OutputMode: OutputYAML,
		OutputPath: "~/.bloop-tunnel/config.yaml",
		HasChanges: false,
		VerificationResults: VerificationResults{
			Connectivity: ConnectivityResult{Success: false, Status: "pending"},
			Enrollment:   EnrollmentResult{Success: false, Status: "pending"},
			Relay:        RelayResult{Success: false, Status: "pending"},
		},
	}
}

func Run(config Config) error {
	p := tea.NewProgram(NewMainModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		return err
	}
	return nil
}
