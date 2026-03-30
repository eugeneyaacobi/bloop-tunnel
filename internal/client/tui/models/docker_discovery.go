package models

import (
	"context"
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"bloop-tunnel/internal/client/dockerdiscover"
)

// DockerDiscoveryMsg is sent when Docker discovery completes
type DockerDiscoveryMsg struct {
	Candidates []dockerdiscover.Candidate
	Error      error
}

// DockerDiscoverySelectMsg is sent when a container is selected
type DockerDiscoverySelectMsg struct {
	Candidate dockerdiscover.Candidate
}

// DockerDiscoverySkipMsg is sent when user skips discovery
type DockerDiscoverySkipMsg struct{}

// DockerDiscoveryBackMsg is sent when user goes back
type DockerDiscoveryBackMsg struct{}

// DockerDiscoveryModel for Docker container discovery
type DockerDiscoveryModel struct {
	loading    bool
	candidates []dockerdiscover.Candidate
	listView   ListViewModel
	errorMsg   string
	width      int
	height     int
	discoverer dockerdiscover.Discoverer
}

// DockerDiscoveryOpts for constructing the model
type DockerDiscoveryOpts struct {
	Discoverer dockerdiscover.Discoverer
}

func NewDockerDiscoveryModel(opts DockerDiscoveryOpts) DockerDiscoveryModel {
	return DockerDiscoveryModel{
		loading:    true,
		candidates: []dockerdiscover.Candidate{},
		listView:   NewListViewModel([]ListItem{}),
		width:      80,
		height:     24,
		discoverer: opts.Discoverer,
	}
}

func (m DockerDiscoveryModel) Init() tea.Cmd {
	if m.discoverer == nil {
		m.discoverer = dockerdiscover.NewClient("")
	}

	return func() tea.Msg {
		ctx := context.Background()
		candidates, err := m.discoverer.Discover(ctx)
		return DockerDiscoveryMsg{Candidates: candidates, Error: err}
	}
}

func (m DockerDiscoveryModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case DockerDiscoveryMsg:
		m.loading = false
		if msg.Error != nil {
			m.errorMsg = fmt.Sprintf("Docker discovery failed: %v", msg.Error)
			return m, nil
		}

		m.candidates = msg.Candidates
		items := make([]ListItem, len(msg.Candidates))
		for i, c := range msg.Candidates {
			items[i] = ListItem{
				ID:      c.ContainerID,
				Label:   c.SuggestedName,
				Details: fmt.Sprintf("  [%s] %s → %s", c.Image, c.ContainerName, c.LocalAddr),
			}
		}
		m.listView = NewListViewModel(items)

	case tea.KeyMsg:
		if m.loading {
			return m, nil
		}

		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "esc":
			return m, func() tea.Msg { return DockerDiscoverySkipMsg{} }
		case "s":
			return m, func() tea.Msg { return DockerDiscoverySkipMsg{} }
		case "b":
			return m, func() tea.Msg { return DockerDiscoveryBackMsg{} }
		case "enter":
			if len(m.candidates) > 0 && m.listView.SelectedIndex() >= 0 {
				idx := m.listView.SelectedIndex()
				return m, func() tea.Msg {
					return DockerDiscoverySelectMsg{Candidate: m.candidates[idx]}
				}
			}
		}

		// Forward navigation to list view
		if !m.loading && len(m.candidates) > 0 {
			var cmd tea.Cmd
			m.listView, cmd = m.listView.Update(msg)
			return m, cmd
		}

	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
		// Update list view dimensions
		visibleHeight := m.height - 12
		if visibleHeight < 5 {
			visibleHeight = 5
		}
		m.listView = *m.listView.WithHeight(visibleHeight).WithWidth(m.width - 8)
	}

	return m, nil
}

func (m DockerDiscoveryModel) View() string {
	titleStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#007AFF")).
		Bold(true).
		MarginTop(1).
		MarginBottom(1)

	headerStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#8E8E93")).
		Italic(true).
		MarginBottom(2)

	footerStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#8E8E93")).
		MarginTop(2)

	if m.loading {
		loadingStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("#007AFF"))

		return lipgloss.JoinVertical(
			lipgloss.Center,
			titleStyle.Render("Docker Service Discovery"),
			"",
			loadingStyle.Render("⠋ Scanning for running Docker containers..."),
		)
	}

	if m.errorMsg != "" {
		errorStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF3B30")).
			MarginBottom(2)

		return lipgloss.JoinVertical(
			lipgloss.Left,
			titleStyle.Render("Docker Service Discovery"),
			"",
			errorStyle.Render(m.errorMsg),
			"",
			footerStyle.Render("Press S to skip • B to go back • Q to quit"),
		)
	}

	if len(m.candidates) == 0 {
		emptyStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("#8E8E93")).
			Italic(true)

		return lipgloss.JoinVertical(
			lipgloss.Left,
			titleStyle.Render("Docker Service Discovery"),
			"",
			headerStyle.Render("No running Docker containers found with exposed HTTP ports."),
			"",
			emptyStyle.Render("You can add tunnels manually in the next step."),
			"",
			footerStyle.Render("Enter: Continue to tunnel setup • S: Skip • Q: Quit"),
		)
	}

	// Has containers - show list
	instructionsStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#8E8E93"))

	return lipgloss.JoinVertical(
		lipgloss.Left,
		titleStyle.Render("Docker Service Discovery"),
		"",
		headerStyle.Render("Select containers to create tunnels for:"),
		"",
		m.listView.View(),
		"",
		instructionsStyle.Render("↑/↓: Navigate • Enter: Add tunnel • S: Skip all • B: Back • Q: Quit"),
	)
}
