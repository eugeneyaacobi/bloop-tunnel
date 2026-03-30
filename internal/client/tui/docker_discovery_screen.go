package tui

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/lipgloss"
)

type DockerContainer struct {
	ID     string
	Name   string
	Image  string
	Ports  []string
	State  string
}

type DockerDiscoveryScreenModel struct {
	config           Config
	width, height    int
	showHelp         bool
	spinner          spinner.Model
	isScanning       bool
	containers       []DockerContainer
	selectedIdx      int
	dockerAvailable  bool
	scanComplete     bool
	skipDocker       bool
}

type DockerScanTickMsg struct{}
type DockerScanCompleteMsg struct {
	Containers []DockerContainer
	Available  bool
}

func NewDockerDiscoveryScreenModel(config Config) DockerDiscoveryScreenModel {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(primaryColor)

	return DockerDiscoveryScreenModel{
		config:          config,
		width:           80,
		height:          24,
		showHelp:        false,
		spinner:         s,
		isScanning:      true,
		containers:      []DockerContainer{},
		selectedIdx:     -1,
		dockerAvailable: false,
		scanComplete:    false,
		skipDocker:      false,
	}
}

func (m DockerDiscoveryScreenModel) Init() tea.Cmd {
	return tea.Batch(
		m.spinner.Tick,
		m.scanDocker(),
	)
}

func (m DockerDiscoveryScreenModel) scanDocker() tea.Cmd {
	return tea.Tick(time.Second*1, func(t time.Time) tea.Msg {
		// Simulate Docker scan (in real implementation, use Docker API)
		// For now, return no containers found
		return DockerScanCompleteMsg{
			Containers: []DockerContainer{},
			Available:  false,
		}
	})
}

func (m DockerDiscoveryScreenModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "enter":
			// Continue to review screen
			return m, func() tea.Msg { return ScreenTransitionMsg{From: StateDockerDiscovery, To: StateReview} }
		case "esc", "s":
			// Skip Docker discovery
			m.skipDocker = true
			return m, func() tea.Msg { return ScreenTransitionMsg{From: StateDockerDiscovery, To: StateReview} }
		case "?":
			m.showHelp = !m.showHelp
			return m, nil
		case "up", "k":
			if m.selectedIdx > 0 {
				m.selectedIdx--
			}
		case "down", "j":
			if m.selectedIdx < len(m.containers)-1 {
				m.selectedIdx++
			}
		case " ":
			// Toggle selection
			if m.selectedIdx >= 0 && m.selectedIdx < len(m.containers) {
				// In a real implementation, this would toggle selection
				// For now, just move to next
			}
		}
	case spinner.TickMsg:
		if m.isScanning {
			var cmd tea.Cmd
			m.spinner, cmd = m.spinner.Update(msg)
			return m, cmd
		}
	case DockerScanCompleteMsg:
		m.isScanning = false
		m.scanComplete = true
		m.containers = msg.Containers
		m.dockerAvailable = msg.Available
		if len(m.containers) > 0 {
			m.selectedIdx = 0
		}
	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
	}

	return m, nil
}

func (m DockerDiscoveryScreenModel) View() string {
	title := titleStyle.Render("Docker Container Discovery")
	subtitle := subtitleStyle.Render("Scan for running containers to auto-create tunnel configurations")

	var content string

	if m.isScanning {
		// Scanning state
		scanSty := lipgloss.NewStyle().
			Foreground(primaryColor).
			MarginTop(2).
			Align(lipgloss.Center)

		content = lipgloss.JoinVertical(
			lipgloss.Center,
			title,
			subtitle,
			"",
			scanSty.Render(m.spinner.View() + " Scanning for Docker containers..."),
			"",
			HelpText("This may take a moment"),
		)
	} else if m.scanComplete && !m.dockerAvailable {
		// Docker not available
		noDockerSty := lipgloss.NewStyle().
			Foreground(mutedColor).
			MarginTop(2).
			Align(lipgloss.Center)

		content = lipgloss.JoinVertical(
			lipgloss.Center,
			title,
			subtitle,
			"",
			noDockerSty.Render("Docker is not available or not running."),
			"",
			noDockerSty.Render("You can configure tunnels manually on the next screen."),
			"",
			buttonStyle.Render("Press Enter to continue"),
		)
	} else if m.scanComplete && len(m.containers) == 0 {
		// No containers found
		noContainersSty := lipgloss.NewStyle().
			Foreground(mutedColor).
			MarginTop(2).
			Align(lipgloss.Center)

		content = lipgloss.JoinVertical(
			lipgloss.Center,
			title,
			subtitle,
			"",
			noContainersSty.Render("No running containers with exposed ports found."),
			"",
			noContainersSty.Render("You can configure tunnels manually on the next screen."),
			"",
			buttonStyle.Render("Press Enter to continue"),
		)
	} else {
		// Show containers
		containerList := m.renderContainers()

		content = lipgloss.JoinVertical(
			lipgloss.Left,
			title,
			subtitle,
			"",
			HelpText("Select containers to create tunnels for, or skip to configure manually"),
			"",
			containerList,
			"",
			buttonStyle.Render("Enter: Continue • Esc: Skip Docker • Space: Toggle selection"),
		)
	}

	// Help footer
	var helpText string
	if m.showHelp {
		helpText = HelpText("Enter: Continue • Esc: Skip • Q: Quit")
	} else {
		helpText = footerStyle.Render("?: Show help")
	}

	fullContent := lipgloss.JoinVertical(lipgloss.Left, content, "", helpText)

	return Center(m.width, fullContent)
}

func (m DockerDiscoveryScreenModel) renderContainers() string {
	if len(m.containers) == 0 {
		return lipgloss.NewStyle().Foreground(mutedColor).Italic(true).Render("No containers found")
	}

	var rows []string
	headerStyle := lipgloss.NewStyle().Foreground(primaryColor).Bold(true)
	rows = append(rows, lipgloss.JoinHorizontal(
		lipgloss.Left,
		headerStyle.Width(25).Render("Container Name"),
		headerStyle.Width(30).Render("Image"),
		headerStyle.Width(15).Render("Ports"),
		headerStyle.Width(10).Render("State"),
	))

	for i, container := range m.containers {
		rowStyle := lipgloss.NewStyle()
		if i == m.selectedIdx {
			rowStyle = rowStyle.Background(mutedColor).Foreground(textColor)
		}

		portsStr := ""
		if len(container.Ports) > 0 {
			portsStr = container.Ports[0]
			if len(container.Ports) > 1 {
				portsStr += " (+more)"
			}
		}

		row := lipgloss.JoinHorizontal(
			lipgloss.Left,
			rowStyle.Width(25).Render(container.Name),
			rowStyle.Width(30).Render(container.Image),
			rowStyle.Width(15).Render(portsStr),
			rowStyle.Width(10).Render(container.State),
		)
		rows = append(rows, row)
	}

	return lipgloss.JoinVertical(lipgloss.Left, rows...)
}

func (m DockerDiscoveryScreenModel) WithContainers(containers []DockerContainer) DockerDiscoveryScreenModel {
	m.containers = containers
	return m
}
