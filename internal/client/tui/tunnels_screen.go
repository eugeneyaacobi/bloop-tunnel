package tui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"bloop-tunnel/internal/client/tui/models"
)

type TunnelsScreenModel struct {
	list         models.ListViewModel
	tunnels      []TunnelConfig
	width, height int
	showHelp     bool
}

func NewTunnelsScreenModel(tunnels []TunnelConfig) TunnelsScreenModel {
	items := make([]models.ListItem, len(tunnels))
	for i, tunnel := range tunnels {
		accessLabel := tunnel.Access
		if tunnel.Access == "basic_auth" {
			accessLabel = "🔒 Basic Auth"
		} else if tunnel.Access == "token_protected" {
			accessLabel = "🔑 Token"
		} else {
			accessLabel = "🌐 Public"
		}

		details := fmt.Sprintf("%s:%d • %s", tunnel.LocalIP, tunnel.LocalPort, accessLabel)
		if tunnel.Hostname != "" {
			details += fmt.Sprintf(" → %s", tunnel.Hostname)
		}

		items[i] = models.ListItem{
			ID:      tunnel.Name,
			Label:   tunnel.Name,
			Details: details,
		}
	}

	list := models.NewListViewModel(items)
	listWidth := 60
	if listWidth > 80 {
		listWidth = 80
	}
	list = *list.WithWidth(listWidth)

	return TunnelsScreenModel{
		list:    list,
		tunnels: tunnels,
		width:   80,
		height:  24,
		showHelp: false,
	}
}

func (m TunnelsScreenModel) Init() tea.Cmd {
	return nil
}

func (m TunnelsScreenModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "a", "n":
			// Add new tunnel
			return m, func() tea.Msg { return TunnelsListAddMsg{} }
		case "e":
			// Edit selected tunnel
			selected := m.list.SelectedIndex()
			return m, func() tea.Msg { return TunnelsListEditMsg{Index: selected} }
		case "d", "delete":
			// Delete selected tunnel
			selected := m.list.SelectedIndex()
			return m, func() tea.Msg { return TunnelDeleteMsg{Index: selected} }
		case "enter":
			// Edit selected tunnel
			selected := m.list.SelectedIndex()
			return m, func() tea.Msg { return TunnelsListEditMsg{Index: selected} }
		case "right", "l":
			// Next screen
			return m, func() tea.Msg { return TunnelsListNextMsg{} }
		case "left", "h", "b":
			// Back to config
			return m, func() tea.Msg { return TunnelsListBackMsg{} }
		case "esc":
			// Back to config
			return m, func() tea.Msg { return TunnelsListBackMsg{} }
		case "?":
			m.showHelp = !m.showHelp
			return m, nil
		case "q", "ctrl+c":
			return m, tea.Quit
		}
	case models.ListSelectMsg:
		// List selection changed - update local state
		m.list = *m.list.MoveSelection(0)
	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
		visibleHeight := m.height - 10
		if visibleHeight < 5 {
			visibleHeight = 5
		}
		m.list = *m.list.WithHeight(visibleHeight)
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m TunnelsScreenModel) View() string {
	title := titleStyle.Render("Tunnels Configuration")

	var subtitle string
	if len(m.tunnels) == 0 {
		subtitle = subtitleStyle.Render("No tunnels configured yet. Add your first tunnel to get started.")
	} else {
		subtitle = subtitleStyle.Render(fmt.Sprintf("You have %d tunnel(s) configured", len(m.tunnels)))
	}

	// Action buttons
	buttonWidth := 12
	addButton := buttonStyle.Width(buttonWidth).Render("Add (a)")
	editButton := buttonStyle.Width(buttonWidth).Render("Edit (e)")
	deleteButton := buttonStyle.Width(buttonWidth).Render("Delete (d)")
	nextButton := buttonStyle.Width(buttonWidth).Render("Next →")
	backButton := buttonStyle.Width(buttonWidth).Render("← Back")

	buttonRow := lipgloss.JoinHorizontal(lipgloss.Left,
		addButton, " ",
		editButton, " ",
		deleteButton, "  ",
		backButton, " ",
		nextButton,
	)

	// List view
	listView := m.list.View()

	// Empty state message
	var listContent string
	if len(m.tunnels) == 0 {
		emptySty := lipgloss.NewStyle().
			Foreground(mutedColor).
			Italic(true).
			MarginTop(2)
		listContent = emptySty.Render("No tunnels configured. Press 'a' to add your first tunnel.")
	} else {
		listContent = listView
	}

	content := lipgloss.JoinVertical(
		lipgloss.Left,
		title,
		subtitle,
		"",
		listContent,
		"",
		buttonRow,
	)

	// Help footer
	var helpText string
	if m.showHelp {
		helpModel := NewHelpModel(ListScreenKeys())
		helpText = helpModel.ShortHelpView()
	} else {
		helpText = footerStyle.Render("?: Show help")
	}

	content = lipgloss.JoinVertical(lipgloss.Left, content, "", helpText)

	return Center(m.width, content)
}

func (m TunnelsScreenModel) WithTunnels(tunnels []TunnelConfig) TunnelsScreenModel {
	m.tunnels = tunnels
	items := make([]models.ListItem, len(tunnels))
	for i, tunnel := range tunnels {
		accessLabel := tunnel.Access
		if tunnel.Access == "basic_auth" {
			accessLabel = "🔒 Basic Auth"
		} else if tunnel.Access == "token_protected" {
			accessLabel = "🔑 Token"
		} else {
			accessLabel = "🌐 Public"
		}

		details := fmt.Sprintf("%s:%d • %s", tunnel.LocalIP, tunnel.LocalPort, accessLabel)
		if tunnel.Hostname != "" {
			details += fmt.Sprintf(" → %s", tunnel.Hostname)
		}

		items[i] = models.ListItem{
			ID:      tunnel.Name,
			Label:   tunnel.Name,
			Details: details,
		}
	}
	m.list = models.NewListViewModel(items)
	return m
}
