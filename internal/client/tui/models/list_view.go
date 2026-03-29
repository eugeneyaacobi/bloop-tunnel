package models

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type ListItem struct {
	ID      string
	Label   string
	Details string // Optional additional info
}

type ListViewModel struct {
	items        []ListItem
	selected     int
	cursorPos    int // Visible scrolling position
	height       int // Number of visible items
	width        int
}

type ListSelectMsg struct{ Index int }
type ListDeleteMsg struct{ Index int }
type ListEditMsg struct{ Index int }

func NewListViewModel(items []ListItem) ListViewModel {
	return ListViewModel{
		items:     items,
		selected:  0,
		cursorPos: 0,
		height:    10, // Default visible items
		width:     60, // Default width
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
		// Reserve space for header/footer
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
		return lipgloss.NewStyle().
			Foreground(lipgloss.Color("#8E8E93")).
			Italic(true).
			Render("No items available.")
	}

	// Calculate visible items
	visibleTop := m.cursorPos
	visibleBottom := min(m.cursorPos+m.height, len(m.items))
	visibleItems := m.items[visibleTop:visibleBottom]

	// Render visible items with selection indicator
	baseStyle := lipgloss.NewStyle().
		Padding(0, 1).
		MarginRight(1).
		Width(m.width - 4)

	selectedStyle := baseStyle.Copy().
		Background(lipgloss.Color("#007AFF")).
		Foreground(lipgloss.Color("#FFFFFF")).
		Bold(true)

	var rows []string
	for i, item := range visibleItems {
		globalIndex := visibleTop + i
		prefix := "  "
		if globalIndex == m.selected {
			prefix = "● "
			rows = append(rows, selectedStyle.Render(prefix+item.Label+item.Details))
		} else {
			rows = append(rows, baseStyle.Render(prefix+item.Label+item.Details))
		}
	}

	// Add scroll indicators
	if visibleTop > 0 {
		rows = append([]string{"▲"}, rows...)
	}
	if visibleBottom < len(m.items)-1 {
		rows = append(rows, "▼")
	}

	boxStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#007AFF")).
		Padding(0, 1).
		Width(m.width)

	return boxStyle.Render(lipgloss.JoinVertical(lipgloss.Left, rows...))
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

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
