package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"bloop-tunnel/internal/client/tui/models"
)

type MainModel struct {
	state          State
	stateStack     []State     // For back navigation
	config         Config
	currentScreen  tea.Model // Current screen model
	previousScreen tea.Model // Previous screen (for back navigation)
	width          int        // Terminal width
	height         int        // Terminal height
}

func NewMainModel() *MainModel {
	return &MainModel{
		state:     StateWelcome,
		stateStack: []State{},
		config:     DefaultConfig(),
		width:      80,
		height:     24,
	}
}

func (m *MainModel) Init() tea.Cmd {
	// Initialize the welcome screen as the current screen
	m.currentScreen = NewWelcomeModel()
	return nil
}

func (m *MainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			if m.config.HasChanges {
				// Show quit confirmation
				m.state = StateQuitConfirm
				m.stateStack = append(m.stateStack, m.state)
				m.currentScreen = NewQuitConfirmModel()
				return m, nil
			}
			return m, tea.Quit
		}
	case ScreenTransitionMsg:
		// Handle screen transitions
		m.state = msg.To
		m.currentScreen = m.getScreenForState(msg.To)
		if msg.From != msg.To {
			m.stateStack = append(m.stateStack, msg.From)
		}
		return m, nil
	case ConfigSaveMsg:
		// Handle config save
		m.config = msg.Config
		m.config.HasChanges = false
		return m, nil
	case ConfigChangedMsg:
		// Mark config as changed
		m.config.HasChanges = true
		return m, nil
	case TunnelAddMsg:
		// Add new tunnel (handled by Endpoints screen)
		m.config.Tunnels = append(m.config.Tunnels, TunnelConfig{})
		m.config.HasChanges = true
		return m, func() tea.Msg { return ScreenTransitionMsg{From: m.state, To: StateTunnels} }
	case TunnelSaveMsg:
		// Save tunnel (update existing)
		if msg.Tunnel.Name != "" {
			found := false
			for i, t := range m.config.Tunnels {
				if t.Name == msg.Tunnel.Name {
					m.config.Tunnels[i] = msg.Tunnel
					found = true
					break
				}
			}
			if !found {
				m.config.Tunnels = append(m.config.Tunnels, msg.Tunnel)
			}
		}
		m.config.HasChanges = true
		return m, func() tea.Msg { return ScreenTransitionMsg{From: m.state, To: StateTunnels} }
	case TunnelDeleteMsg:
		// Delete tunnel
		if msg.Index >= 0 && msg.Index < len(m.config.Tunnels) {
			m.config.Tunnels = append(m.config.Tunnels[:msg.Index], m.config.Tunnels[msg.Index+1:]...)
		}
		m.config.HasChanges = true
		return m, nil
	case TunnelsListAddMsg:
		// Navigate to Endpoints screen for new tunnel
		return m, func() tea.Msg { return ScreenTransitionMsg{From: StateTunnels, To: StateEndpoints} }
	case TunnelsListEditMsg:
		// Navigate to Endpoints screen with tunnel data
		if msg.Index >= 0 && msg.Index < len(m.config.Tunnels) {
			// TODO: Handled by main router, need to pass tunnel data
			// For now, just show a stub screen
			m.state = StateEndpoints
			m.stateStack = append(m.stateStack, StateTunnels)
			m.currentScreen = models.NewStatusModel().WithSuccess("Tunnel form not yet implemented")
			return m, nil
		}
		return m, nil
	case TunnelsListNextMsg:
		// Navigate to verification/review
		return m, func() tea.Msg { return ScreenTransitionMsg{From: StateTunnels, To: StateReview} }
	case TunnelsListBackMsg:
		// Navigate back to config
		return m, func() tea.Msg { return ScreenTransitionMsg{From: StateTunnels, To: StateConfig} }
	case QuitConfirmMsg:
		if msg.Choice == "Quit" {
			return m, tea.Quit
		}
		// Cancel - go back to previous state
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

	// Delegate update to current screen
	if m.currentScreen != nil {
		var cmd tea.Cmd
		m.currentScreen, cmd = m.currentScreen.Update(msg)
		return m, cmd
	}

	return m, nil
}

func (m MainModel) getScreenForState(state State) tea.Model {
	switch state {
	case StateWelcome:
		return NewWelcomeModel()
	case StateConfig:
		// TODO: Implement config screen - needs to be moved from screens package
		return models.NewStatusModel().WithSuccess("Config screen not yet implemented")
	case StateEndpoints:
		// TODO: Implement tunnel form - needs to be moved from screens package
		return models.NewStatusModel().WithSuccess("Tunnel form not yet implemented")
	case StateTunnels:
		// TODO: Implement tunnels list - needs to be moved from screens package
		return models.NewStatusModel().WithSuccess("Tunnels list not yet implemented")
	case StateVerification:
		// TODO: Implement verification screen
		return models.NewStatusModel()
	case StateDockerDiscovery:
		// TODO: Implement Docker discovery screen
		return models.NewStatusModel()
	case StateReview:
		// TODO: Implement review screen
		return models.NewStatusModel()
	case StateOutput:
		// TODO: Implement output screen
		return models.NewStatusModel()
	case StateQuitConfirm:
		return NewQuitConfirmModel()
	default:
		return models.NewStatusModel()
	}
}

func (m *MainModel) View() string {
	// Delegate view to current screen
	if m.currentScreen != nil {
		return m.currentScreen.View()
	}

	// Fallback
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FF3B30")).
		Render("Error: No screen active")
}

// Default config for new setup
func DefaultConfig() Config {
	return Config{
		ControlPlaneURL:      "https://api.bloop.to",
		RelayURL:             "wss://relay.bloop.to/connect",
		AuthTokenEnv:         "BLOOP_CLIENT_TOKEN",
		EnrollmentTokenEnv:   "",
		Reconnect: ReconnectConfig{
			InitialDelayMs: 1000,
			MaxDelayMs:     30000,
		},
		Tunnels:              []TunnelConfig{},
		OutputMode:           OutputYAML,
		OutputPath:           "~/.bloop-tunnel/config.yaml",
		HasChanges:           false,
		VerificationResults: VerificationResults{
			Connectivity: ConnectivityResult{Success: false, Status: "pending"},
			Enrollment:  EnrollmentResult{Success: false, Status: "pending"},
			Relay:       RelayResult{Success: false, Status: "pending"},
		},
	}
}

// Run the TUI
func Run(config Config) error {
	p := tea.NewProgram(NewMainModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		return err
	}
	return nil
}
