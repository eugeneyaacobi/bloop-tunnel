package tui

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/help"
	tea "github.com/charmbracelet/bubbletea"
)

// KeyMap defines the keybindings used across all screens
type KeyMap struct {
	Up     key.Binding
	Down   key.Binding
	Left   key.Binding
	Right  key.Binding
	Enter  key.Binding
	Back   key.Binding
	Quit   key.Binding
	Help   key.Binding
	Delete key.Binding
	Add    key.Binding
	Edit   key.Binding
	Save   key.Binding
	Cancel key.Binding
	Next   key.Binding
	Prev   key.Binding
	Tab    key.Binding
}

// Default keymap for the application
var Keys = KeyMap{
	Up: key.NewBinding(
		key.WithKeys("up", "k"),
		key.WithHelp("↑/k", "move up"),
	),
	Down: key.NewBinding(
		key.WithKeys("down", "j"),
		key.WithHelp("↓/j", "move down"),
	),
	Left: key.NewBinding(
		key.WithKeys("left", "h"),
		key.WithHelp("←/h", "move left"),
	),
	Right: key.NewBinding(
		key.WithKeys("right", "l"),
		key.WithHelp("→/l", "move right"),
	),
	Enter: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "select"),
	),
	Back: key.NewBinding(
		key.WithKeys("esc", "b"),
		key.WithHelp("esc/b", "back"),
	),
	Quit: key.NewBinding(
		key.WithKeys("q", "ctrl+c"),
		key.WithHelp("q", "quit"),
	),
	Help: key.NewBinding(
		key.WithKeys("?"),
		key.WithHelp("?", "toggle help"),
	),
	Delete: key.NewBinding(
		key.WithKeys("d", "delete"),
		key.WithHelp("d", "delete"),
	),
	Add: key.NewBinding(
		key.WithKeys("a", "n"),
		key.WithHelp("a/n", "add new"),
	),
	Edit: key.NewBinding(
		key.WithKeys("e"),
		key.WithHelp("e", "edit"),
	),
	Save: key.NewBinding(
		key.WithKeys("ctrl+s"),
		key.WithHelp("ctrl+s", "save"),
	),
	Cancel: key.NewBinding(
		key.WithKeys("ctrl+c", "esc"),
		key.WithHelp("esc", "cancel"),
	),
	Next: key.NewBinding(
		key.WithKeys("n", "right"),
		key.WithHelp("n/→", "next"),
	),
	Prev: key.NewBinding(
		key.WithKeys("p", "left"),
		key.WithHelp("p/←", "previous"),
	),
	Tab: key.NewBinding(
		key.WithKeys("tab"),
		key.WithHelp("tab", "next field"),
	),
}

// bindingsKeyMap adapts []key.Binding to help.KeyMap
type bindingsKeyMap []key.Binding

func (b bindingsKeyMap) ShortHelp() []key.Binding { return b }
func (b bindingsKeyMap) FullHelp() [][]key.Binding { return [][]key.Binding{b} }

// HelpModel wraps the bubbles/help component
type HelpModel struct {
	help help.Model
	keys bindingsKeyMap
}

// NewHelpModel creates a new help model
func NewHelpModel(bindings []key.Binding) HelpModel {
	h := help.New()
	h.ShowAll = true
	return HelpModel{
		help: h,
		keys: bindingsKeyMap(bindings),
	}
}

func (m HelpModel) Init() tea.Cmd {
	return nil
}

func (m HelpModel) Update(msg tea.Msg) (HelpModel, tea.Cmd) {
	var cmd tea.Cmd
	m.help, cmd = m.help.Update(msg)
	return m, cmd
}

func (m HelpModel) View() string {
	return m.help.View(m.keys)
}

func (m HelpModel) ShortHelpView() string {
	return m.help.ShortHelpView(m.keys)
}

// Screen-specific help key groups

// ConfigScreenKeys returns keybindings for the config screen
func ConfigScreenKeys() []key.Binding {
	return []key.Binding{Keys.Up, Keys.Down, Keys.Enter, Keys.Tab, Keys.Save, Keys.Quit}
}

// FormScreenKeys returns keybindings for form screens (endpoints)
func FormScreenKeys() []key.Binding {
	return []key.Binding{Keys.Up, Keys.Down, Keys.Tab, Keys.Enter, Keys.Save, Keys.Cancel}
}

// ListScreenKeys returns keybindings for list screens (tunnels)
func ListScreenKeys() []key.Binding {
	return []key.Binding{Keys.Up, Keys.Down, Keys.Add, Keys.Edit, Keys.Delete, Keys.Next, Keys.Back}
}

// VerificationScreenKeys returns keybindings for verification screen
func VerificationScreenKeys() []key.Binding {
	return []key.Binding{Keys.Up, Keys.Down, Keys.Enter, Keys.Quit}
}

// ReviewScreenKeys returns keybindings for review screen
func ReviewScreenKeys() []key.Binding {
	return []key.Binding{Keys.Up, Keys.Down, Keys.Back, Keys.Enter}
}

// OutputScreenKeys returns keybindings for output screen
func OutputScreenKeys() []key.Binding {
	return []key.Binding{Keys.Quit, Keys.Back}
}
