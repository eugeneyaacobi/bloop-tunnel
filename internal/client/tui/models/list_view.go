package models

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type ListItem struct {
	ID      string
	Label   string
	Details string // Optional additional info
}

type ListViewModel struct {
	items     []ListItem
	selected  int
	cursorPos int // Visible scrolling position
	height    int // Number of visible items
	width     int
}

type ListSelectMsg struct{ Index int }
type ListDeleteMsg struct{ Index int }
type ListEditMsg struct{ Index int }

var (
	listEmptySty    = lipgloss.NewStyle().Foreground(lipgloss.Color("#8E8E93")).Italic(true)
	listItemSty     = lipgloss.NewStyle().Padding(0, 1).MarginRight(1)
	listSelectedSty = lipgloss.NewStyle().Padding(0, 1).MarginRight(1).Background(lipgloss.Color("#007AFF")).Foreground(lipgloss.Color("#FFFFFF")).Bold(true)
	listBoxSty      = lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color("#007AFF")).Padding(0, 1)
)

func NewListViewModel(items []ListItem) ListViewModel {
	return ListViewModel{
		items:     items,
		selected:  0,
		cursorPos: 0,
		height:    10,
		width:     60,
	}
}

func (m ListViewModel) Init() tea.Cmd {
	return nil
}

func (m ListViewModel) Update(msg tea.Msg) (ListViewModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			if m.selected > 0 {
				m.selected--
				m.scrollUp()
			}
		case "down", "j":
			if m.selected < len(m.items)-1 {
				m.selected++
				m.scrollDown()
			}
		case "home", "g":
			m.selected = 0
			m.cursorPos = 0
		case "end", "G":
			m.selected = len(m.items) - 1
			m.scrollToEnd()
		case "pageup", "ctrl+b":
			m.selected = max(0, m.selected-m.height)
			m.scrollUp()
		case "pagedown", "ctrl+f":
			m.selected = min(len(m.items)-1, m.selected+m.height)
			m.scrollDown()
		case "enter", " ":
			return m, func() tea.Msg { return ListSelectMsg{Index: m.selected} }
		}
	case ListSelectMsg:
		m.selected = msg.Index
		m.scrollToShowSelection()
	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
		visibleHeight := m.height - 10
		if visibleHeight < 5 {
			visibleHeight = 5
		}
		m.height = visibleHeight
		m.scrollToShowSelection()
	}

	return m, nil
}

func (m *ListViewModel) scrollUp() {
	if m.selected < m.cursorPos {
		m.cursorPos = m.selected
	}
}

func (m *ListViewModel) scrollDown() {
	visibleBottom := m.cursorPos + m.height - 1
	if m.selected > visibleBottom {
		m.cursorPos = m.selected - m.height + 1
	}
}

func (m *ListViewModel) scrollToShowSelection() {
	if m.selected < m.cursorPos {
		m.cursorPos = m.selected
	} else if m.selected > m.cursorPos+m.height-1 {
		m.cursorPos = m.selected - m.height + 1
	}
}

func (m *ListViewModel) scrollToEnd() {
	if len(m.items) > m.height {
		m.cursorPos = len(m.items) - m.height
	}
}

func (m *ListViewModel) WithHeight(height int) *ListViewModel {
	m.height = height
	if m.height > len(m.items) {
		m.height = len(m.items)
	}
	return m
}

func (m *ListViewModel) WithWidth(width int) *ListViewModel {
	m.width = width
	return m
}

func (m *ListViewModel) MoveSelection(delta int) *ListViewModel {
	m.selected += delta
	if m.selected < 0 {
		m.selected = 0
	} else if m.selected >= len(m.items) {
		m.selected = len(m.items) - 1
	}
	m.scrollToShowSelection()
	return m
}

func (m *ListViewModel) View() string {
	if len(m.items) == 0 {
		return listEmptySty.Render("No items available.")
	}

	visibleTop := m.cursorPos
	visibleBottom := min(m.cursorPos+m.height, len(m.items))
	visibleItems := m.items[visibleTop:visibleBottom]

	itemWidth := m.width - 4
	if itemWidth < 10 {
		itemWidth = 10
	}

	baseSty := listItemSty.Width(itemWidth)
	selSty := listSelectedSty.Width(itemWidth)

	var rows []string
	for i, item := range visibleItems {
		globalIndex := visibleTop + i
		text := item.Label
		if item.Details != "" {
			text = fmt.Sprintf("%s %s", item.Label, item.Details)
		}
		if globalIndex == m.selected {
			rows = append(rows, selSty.Render("● "+text))
		} else {
			rows = append(rows, baseSty.Render("  "+text))
		}
	}

	if visibleTop > 0 {
		rows = append([]string{"▲"}, rows...)
	}
	if visibleBottom < len(m.items) {
		rows = append(rows, "▼")
	}

	return listBoxSty.Width(m.width).Render(lipgloss.JoinVertical(lipgloss.Left, rows...))
}

func (m ListViewModel) SelectedIndex() int {
	return m.selected
}

func (m ListViewModel) Selected() ListItem {
	if m.selected >= 0 && m.selected < len(m.items) {
		return m.items[m.selected]
	}
	return ListItem{}
}
