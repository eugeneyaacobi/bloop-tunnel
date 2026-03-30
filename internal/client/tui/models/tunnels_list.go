package models

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// TunnelsListAddMsg is sent when user wants to add a new tunnel
type TunnelsListAddMsg struct{}

// TunnelsListEditMsg is sent when user wants to edit a tunnel
type TunnelsListEditMsg struct {
	Index int
}

// TunnelsListDeleteMsg is sent when user wants to delete a tunnel
type TunnelsListDeleteMsg struct {
	Index int
}

// TunnelsListNextMsg is sent when user wants to proceed
type TunnelsListNextMsg struct{}

// TunnelsListBackMsg is sent when user wants to go back
type TunnelsListBackMsg struct{}

// TunnelsListModel manages the list of configured tunnels
type TunnelsListModel struct {
	tunnels   []TunnelConfigEntry
	listView  ListViewModel
	width     int
	height    int
	helpText  string
}

// TunnelConfigEntry represents a tunnel in the list
type TunnelConfigEntry struct {
	Index   int
	Name    string
	Address string
	Access  string
	Details string
}

// TunnelsListOpts for constructing the model
type TunnelsListOpts struct {
	Tunnels []TunnelConfigEntry
}

func NewTunnelsListModel(opts TunnelsListOpts) TunnelsListModel {
	items := make([]ListItem, len(opts.Tunnels))
	for i, t := range opts.Tunnels {
		items[i] = ListItem{
			ID:      t.Name,
			Label:   t.Name,
			Details: "  " + t.Address + " [" + t.Access + "]",
		}
	}

	return TunnelsListModel{
		tunnels:  opts.Tunnels,
		listView: NewListViewModel(items),
		width:    80,
		height:   24,
		helpText: "↑/↓: Navigate • A: Add • E: Edit • D: Delete • Enter: Next • B: Back • Q: Quit",
	}
}

func (m TunnelsListModel) Init() tea.Cmd {
	return nil
}

func (m TunnelsListModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "a":
			return m, func() tea.Msg { return TunnelsListAddMsg{} }
		case "e":
			if m.listView.SelectedIndex() >= 0 && len(m.tunnels) > 0 {
				return m, func() tea.Msg {
					return TunnelsListEditMsg{Index: m.listView.SelectedIndex()}
				}
			}
		case "d":
			if m.listView.SelectedIndex() >= 0 && len(m.tunnels) > 0 {
				return m, func() tea.Msg {
					return TunnelsListDeleteMsg{Index: m.listView.SelectedIndex()}
				}
			}
		case "enter":
			return m, func() tea.Msg { return TunnelsListNextMsg{} }
		case "b", "esc":
			return m, func() tea.Msg { return TunnelsListBackMsg{} }
		}

		// Forward navigation to list view
		var cmd tea.Cmd
		m.listView, cmd = m.listView.Update(msg)
		return m, cmd

	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
		visibleHeight := m.height - 10
		if visibleHeight < 5 {
			visibleHeight = 5
		}
		m.listView = *m.listView.WithHeight(visibleHeight).WithWidth(m.width - 8)
	}

	return m, nil
}

func (m TunnelsListModel) View() string {
	titleStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#007AFF")).
		Bold(true).
		MarginTop(1).
		MarginBottom(1)

	headerStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#8E8E93")).
		Italic(true).
		MarginBottom(1)

	countStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#8E8E93")).
		MarginBottom(2)

	footerStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#8E8E93")).
		MarginTop(2)

	emptyStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#8E8E93")).
		Italic(true)

	var content string

	if len(m.tunnels) == 0 {
		content = emptyStyle.Render("No tunnels configured yet. Press 'A' to add your first tunnel.")
	} else {
		content = m.listView.View()
	}

	return lipgloss.JoinVertical(
		lipgloss.Left,
		titleStyle.Render("Tunnel Configuration"),
		"",
		headerStyle.Render("Manage your tunnel endpoints:"),
		countStyle.Render(fmt.Sprintf("Configured tunnels: %d", len(m.tunnels))),
		content,
		"",
		footerStyle.Render(m.helpText),
	)
}

// WithTunnels updates the tunnels list
func (m TunnelsListModel) WithTunnels(tunnels []TunnelConfigEntry) TunnelsListModel {
	m.tunnels = tunnels
	items := make([]ListItem, len(tunnels))
	for i, t := range tunnels {
		items[i] = ListItem{
			ID:      t.Name,
			Label:   t.Name,
			Details: "  " + t.Address + " [" + t.Access + "]",
		}
	}
	m.listView = NewListViewModel(items)
	return m
}

// SelectedIndex returns the currently selected index
func (m TunnelsListModel) SelectedIndex() int {
	return m.listView.SelectedIndex()
}
