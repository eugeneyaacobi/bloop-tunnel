package tui

import (
	"fmt"
	"net"
	"strconv"

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
		return m, func() tea.Msg { return ScreenTransitionMsg{From: m.state, To: StateEndpoints} }
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
	case models.TunnelsListAddMsg:
		// Navigate to tunnel form for new tunnel
		m.state = StateEndpoints
		m.stateStack = append(m.stateStack, StateTunnels)
		m.currentScreen = models.NewTunnelFormModel(models.TunnelFormOpts{
			Tunnel:    models.TunnelFormData{},
			IsEditing: false,
		})
		return m, nil
	case models.TunnelsListEditMsg:
		// Navigate to tunnel form with existing tunnel data
		if msg.Index >= 0 && msg.Index < len(m.config.Tunnels) {
			tunnel := m.config.Tunnels[msg.Index]
			tunnelData := models.TunnelFormData{
				Name:      tunnel.Name,
				Hostname:  tunnel.Hostname,
				LocalIP:   tunnel.LocalIP,
				LocalPort: fmt.Sprintf("%d", tunnel.LocalPort),
				Access:    tunnel.Access,
				TokenEnv:  tunnel.TokenEnv,
			}
			if tunnel.Access == "basic_auth" {
				tunnelData.BasicAuth = models.BasicAuthFormData{
					Username:    tunnel.BasicAuth.Username,
					PasswordEnv: tunnel.BasicAuth.PasswordEnv,
				}
			}
			m.state = StateEndpoints
			m.stateStack = append(m.stateStack, StateTunnels)
			m.currentScreen = models.NewTunnelFormModel(models.TunnelFormOpts{
				Tunnel:    tunnelData,
				IsEditing: true,
			})
			return m, nil
		}
		return m, nil
	case models.TunnelsListDeleteMsg:
		// Delete tunnel
		if msg.Index >= 0 && msg.Index < len(m.config.Tunnels) {
			m.config.Tunnels = append(m.config.Tunnels[:msg.Index], m.config.Tunnels[msg.Index+1:]...)
		}
		m.config.HasChanges = true
		return m, nil
	case models.TunnelsListNextMsg:
		// Navigate to verification/review
		// Build review data
		globalConfig := models.GlobalConfigReview{
			ControlPlaneURL: m.config.ControlPlaneURL,
			RelayURL:        m.config.RelayURL,
			AuthTokenEnv:    m.config.AuthTokenEnv,
		}
		tunnels := make([]models.TunnelConfigReview, len(m.config.Tunnels))
		for i, t := range m.config.Tunnels {
			authDetails := ""
			if t.Access == "basic_auth" {
				authDetails = fmt.Sprintf("user: %s, pass env: %s", t.BasicAuth.Username, t.BasicAuth.PasswordEnv)
			} else if t.Access == "token_protected" {
				authDetails = fmt.Sprintf("token env: %s", t.TokenEnv)
			}
			localAddr := fmt.Sprintf("%s:%d", t.LocalIP, t.LocalPort)
			tunnels[i] = models.TunnelConfigReview{
				Name:        t.Name,
				Hostname:    t.Hostname,
				LocalAddr:   localAddr,
				Access:      t.Access,
				AuthDetails: authDetails,
			}
		}
		reviewData := models.ReviewScreenData{
			Global:     globalConfig,
			Tunnels:    tunnels,
			Output:     string(m.config.OutputMode),
			OutputPath: m.config.OutputPath,
		}
		m.state = StateReview
		m.stateStack = append(m.stateStack, StateTunnels)
		m.currentScreen = models.NewReviewScreenModel(models.ReviewScreenOpts{Data: reviewData})
		return m, nil
	case models.TunnelsListBackMsg:
		// Navigate back to config
		return m, func() tea.Msg { return ScreenTransitionMsg{From: StateTunnels, To: StateConfig} }
	case models.TunnelFormSaveMsg:
		// Save tunnel form data
		port, _ := strconv.Atoi(msg.Tunnel.LocalPort)
		tunnel := TunnelConfig{
			Name:      msg.Tunnel.Name,
			Hostname:  msg.Tunnel.Hostname,
			LocalIP:   msg.Tunnel.LocalIP,
			LocalPort: port,
			Access:    msg.Tunnel.Access,
			TokenEnv:  msg.Tunnel.TokenEnv,
		}
		if msg.Tunnel.Access == "basic_auth" {
			tunnel.BasicAuth = BasicAuthConfig{
				Username:    msg.Tunnel.BasicAuth.Username,
				PasswordEnv: msg.Tunnel.BasicAuth.PasswordEnv,
			}
		}
		// Update or add tunnel
		found := false
		for i, t := range m.config.Tunnels {
			if t.Name == tunnel.Name {
				m.config.Tunnels[i] = tunnel
				found = true
				break
			}
		}
		if !found {
			m.config.Tunnels = append(m.config.Tunnels, tunnel)
		}
		m.config.HasChanges = true
		return m, func() tea.Msg { return ScreenTransitionMsg{From: StateEndpoints, To: StateTunnels} }
	case models.TunnelFormCancelMsg:
		// Cancel - go back to tunnels list
		return m, func() tea.Msg { return ScreenTransitionMsg{From: StateEndpoints, To: StateTunnels} }
	case models.ConfigScreenSaveMsg, models.ConfigScreenSkipMsg:
		// Save config settings
		if configMsg, ok := msg.(models.ConfigScreenSaveMsg); ok {
			m.config.ControlPlaneURL = configMsg.Config.ControlPlaneURL
			m.config.RelayURL = configMsg.Config.RelayURL
			m.config.AuthTokenEnv = configMsg.Config.AuthTokenEnv
		}
		// Navigate to tunnels list
		tunnelEntries := make([]models.TunnelConfigEntry, len(m.config.Tunnels))
		for i, t := range m.config.Tunnels {
			authDetails := ""
			if t.Access == "basic_auth" {
				authDetails = fmt.Sprintf("user: %s", t.BasicAuth.Username)
			} else if t.Access == "token_protected" {
				authDetails = "token"
			}
			localAddr := fmt.Sprintf("%s:%d", t.LocalIP, t.LocalPort)
			tunnelEntries[i] = models.TunnelConfigEntry{
				Index:   i,
				Name:    t.Name,
				Address: localAddr,
				Access:  t.Access,
				Details: authDetails,
			}
		}
		m.state = StateTunnels
		m.stateStack = append(m.stateStack, StateConfig)
		m.currentScreen = models.NewTunnelsListModel(models.TunnelsListOpts{Tunnels: tunnelEntries})
		return m, nil
	case models.ConfigScreenCancelMsg:
		// Cancel - go back to welcome
		return m, func() tea.Msg { return ScreenTransitionMsg{From: StateConfig, To: StateWelcome} }
	case models.DockerDiscoverySelectMsg:
		// Add tunnel from Docker discovery
		candidate := msg.Candidate
		// Parse local address to get IP and port
		localIP, localPortStr, err := parseLocalAddr(candidate.LocalAddr)
		localPort := 8080
		if err == nil {
			port, _ := strconv.Atoi(localPortStr)
			if port > 0 {
				localPort = port
			}
		}
		tunnel := TunnelConfig{
			Name:      candidate.SuggestedName,
			LocalIP:   localIP,
			LocalPort: localPort,
			Access:    "public",
		}
		// Add tunnel
		found := false
		for i, t := range m.config.Tunnels {
			if t.Name == tunnel.Name {
				m.config.Tunnels[i] = tunnel
				found = true
				break
			}
		}
		if !found {
			m.config.Tunnels = append(m.config.Tunnels, tunnel)
		}
		m.config.HasChanges = true
		// Go to tunnels list
		tunnelEntries := make([]models.TunnelConfigEntry, len(m.config.Tunnels))
		for i, t := range m.config.Tunnels {
			localAddr := fmt.Sprintf("%s:%d", t.LocalIP, t.LocalPort)
			tunnelEntries[i] = models.TunnelConfigEntry{
				Index:   i,
				Name:    t.Name,
				Address: localAddr,
				Access:  t.Access,
			}
		}
		m.state = StateTunnels
		m.stateStack = append(m.stateStack, StateDockerDiscovery)
		m.currentScreen = models.NewTunnelsListModel(models.TunnelsListOpts{Tunnels: tunnelEntries})
		return m, nil
	case models.DockerDiscoverySkipMsg, models.DockerDiscoveryBackMsg:
		// Skip or back - go to tunnels list
		tunnelEntries := make([]models.TunnelConfigEntry, len(m.config.Tunnels))
		for i, t := range m.config.Tunnels {
			localAddr := fmt.Sprintf("%s:%d", t.LocalIP, t.LocalPort)
			tunnelEntries[i] = models.TunnelConfigEntry{
				Index:   i,
				Name:    t.Name,
				Address: localAddr,
				Access:  t.Access,
			}
		}
		m.state = StateTunnels
		m.currentScreen = models.NewTunnelsListModel(models.TunnelsListOpts{Tunnels: tunnelEntries})
		return m, nil
	case models.ReviewScreenConfirmMsg:
		// Confirm and proceed to output
		m.state = StateOutput
		m.stateStack = append(m.stateStack, StateReview)
		m.currentScreen = models.NewOutputScreenModel(models.OutputScreenOpts{
			OutputMode: string(m.config.OutputMode),
			OutputPath: m.config.OutputPath,
		})
		return m, nil
	case models.ReviewScreenBackMsg:
		// Go back to tunnels list
		return m, func() tea.Msg { return ScreenTransitionMsg{From: StateReview, To: StateTunnels} }
	case models.ReviewScreenEditMsg:
		// Edit config section
		if msg.Section == "global" {
			return m, func() tea.Msg { return ScreenTransitionMsg{From: StateReview, To: StateConfig} }
		}
		return m, func() tea.Msg { return ScreenTransitionMsg{From: StateReview, To: StateTunnels} }
	case models.OutputScreenConfirmMsg:
		// Finalize output settings
		m.config.OutputMode = OutputMode(msg.Mode)
		m.config.OutputPath = msg.Path
		m.config.HasChanges = false
		// Here you would write the actual output file
		return m, tea.Quit
	case models.OutputScreenBackMsg:
		// Go back to review
		return m, func() tea.Msg { return ScreenTransitionMsg{From: StateOutput, To: StateReview} }
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
		configData := models.ConfigScreenData{
			ControlPlaneURL: m.config.ControlPlaneURL,
			RelayURL:        m.config.RelayURL,
			AuthTokenEnv:    m.config.AuthTokenEnv,
		}
		return models.NewConfigScreenModel(models.ConfigScreenOpts{
			Config:       configData,
			ShowAdvanced: false,
		})
	case StateEndpoints:
		// Tunnel form - should be initialized with data before showing
		return models.NewTunnelFormModel(models.TunnelFormOpts{
			Tunnel:    models.TunnelFormData{},
			IsEditing: false,
		})
	case StateTunnels:
		// Tunnels list - should be initialized with data before showing
		tunnelEntries := make([]models.TunnelConfigEntry, len(m.config.Tunnels))
		for i, t := range m.config.Tunnels {
			localAddr := fmt.Sprintf("%s:%d", t.LocalIP, t.LocalPort)
			tunnelEntries[i] = models.TunnelConfigEntry{
				Index:   i,
				Name:    t.Name,
				Address: localAddr,
				Access:  t.Access,
			}
		}
		return models.NewTunnelsListModel(models.TunnelsListOpts{Tunnels: tunnelEntries})
	case StateVerification:
		return models.NewStatusModel()
	case StateDockerDiscovery:
		return models.NewDockerDiscoveryModel(models.DockerDiscoveryOpts{
			Discoverer: nil,
		})
	case StateReview:
		// Review screen - should be initialized with data before showing
		globalConfig := models.GlobalConfigReview{
			ControlPlaneURL: m.config.ControlPlaneURL,
			RelayURL:        m.config.RelayURL,
			AuthTokenEnv:    m.config.AuthTokenEnv,
		}
		tunnels := make([]models.TunnelConfigReview, len(m.config.Tunnels))
		for i, t := range m.config.Tunnels {
			authDetails := ""
			if t.Access == "basic_auth" {
				authDetails = fmt.Sprintf("user: %s", t.BasicAuth.Username)
			} else if t.Access == "token_protected" {
				authDetails = "token"
			}
			localAddr := fmt.Sprintf("%s:%d", t.LocalIP, t.LocalPort)
			tunnels[i] = models.TunnelConfigReview{
				Name:        t.Name,
				Hostname:    t.Hostname,
				LocalAddr:   localAddr,
				Access:      t.Access,
				AuthDetails: authDetails,
			}
		}
		reviewData := models.ReviewScreenData{
			Global:     globalConfig,
			Tunnels:    tunnels,
			Output:     string(m.config.OutputMode),
			OutputPath: m.config.OutputPath,
		}
		return models.NewReviewScreenModel(models.ReviewScreenOpts{Data: reviewData})
	case StateOutput:
		return models.NewOutputScreenModel(models.OutputScreenOpts{
			OutputMode: string(m.config.OutputMode),
			OutputPath: m.config.OutputPath,
		})
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

// parseLocalAddr parses "host:port" into components
func parseLocalAddr(addr string) (host, port string, err error) {
	h, p, err := net.SplitHostPort(addr)
	if err != nil {
		return "", "", err
	}
	return h, p, nil
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
