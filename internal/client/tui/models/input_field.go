package models

import (
	"strings"
	"unicode/utf8"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type InputFieldModel struct {
	label       string
	placeholder string
	value       string
	focused     bool
	isPassword  bool
	validation  func(string) error
	errorMsg    string
	cursorPos   int
}

type InputFieldOpts struct {
	Label       string
	Placeholder string
	Value       string
	Validation  func(string) error
	IsPassword  bool
}

type InputFieldFocusedMsg struct{ Index int }
type InputFieldValueChangeMsg struct {
	Index int
	Value string
}

func NewInputField(opts InputFieldOpts) InputFieldModel {
	cursorPos := len(opts.Value)
	return InputFieldModel{
		label:       opts.Label,
		placeholder: opts.Placeholder,
		value:       opts.Value,
		focused:     false,
		isPassword:  opts.IsPassword,
		validation:  opts.Validation,
		cursorPos:   cursorPos,
	}
}

func (m InputFieldModel) Init() tea.Cmd {
	return nil
}

func (m InputFieldModel) Update(msg tea.Msg) (InputFieldModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "left":
			if m.cursorPos > 0 {
				m.cursorPos--
			}
		case "right":
			if m.cursorPos < utf8.RuneCountInString(m.value) {
				m.cursorPos++
			}
		case "backspace":
			if m.cursorPos > 0 {
				runes := []rune(m.value)
				m.value = string(append(runes[:m.cursorPos-1], runes[m.cursorPos:]...))
				m.cursorPos--
				m.validate()
			}
		case "delete":
			if m.cursorPos < utf8.RuneCountInString(m.value) {
				runes := []rune(m.value)
				m.value = string(append(runes[:m.cursorPos], runes[m.cursorPos+1:]...))
				m.validate()
			}
		case "ctrl+a":
			m.cursorPos = 0
		case "ctrl+e":
			m.cursorPos = utf8.RuneCountInString(m.value)
		case "enter":
			// Handled by parent model
		}
	case InputFieldValueChangeMsg:
		if msg.Index >= 0 {
			m.value = msg.Value
			m.cursorPos = utf8.RuneCountInString(m.value)
			m.validate()
		}
	case InputFieldFocusedMsg:
		m.focused = true
	default:
		// Handle rune input for typing
		if keyMsg, ok := msg.(tea.KeyMsg); ok {
			if keyMsg.Type == tea.KeyRunes {
				runes := []rune(m.value)
				m.value = string(append(runes[:m.cursorPos], append(keyMsg.Runes, runes[m.cursorPos:]...)...))
				m.cursorPos += len(keyMsg.Runes)
				m.validate()
			}
		}
	}

	return m, nil
}

func (m *InputFieldModel) validate() {
	if m.validation != nil {
		if err := m.validation(m.value); err != nil {
			m.errorMsg = err.Error()
		} else {
			m.errorMsg = ""
		}
	}
}

func (m InputFieldModel) View() string {
	labelSty := lipgloss.NewStyle().Foreground(lipgloss.Color("#007AFF")).Bold(true)
	valueSty := lipgloss.NewStyle().Foreground(lipgloss.Color("#FFFFFF"))
	errSty := lipgloss.NewStyle().Foreground(lipgloss.Color("#FF3B30"))

	// Empty state — show placeholder
	if m.value == "" && m.placeholder != "" && !m.isPassword {
		placeholderSty := lipgloss.NewStyle().Foreground(lipgloss.Color("#8E8E93")).Italic(true)
		return labelSty.Render(m.label+":") + "\n" + placeholderSty.Render(m.placeholder)
	}

	// Build display value with cursor
	runes := []rune(m.value)
	if m.isPassword {
		runes = []rune(strings.Repeat("•", utf8.RuneCountInString(m.value)))
	}

	// Insert cursor block at position
	var display string
	if m.focused {
		before := string(runes[:m.cursorPos])
		after := ""
		if m.cursorPos < len(runes) {
			after = string(runes[m.cursorPos:])
		}
		display = before + "▊" + after
	} else {
		display = string(runes)
	}

	var b strings.Builder
	b.WriteString(labelSty.Render(m.label + ":"))
	b.WriteString("\n")
	b.WriteString(valueSty.Render(display))

	if m.errorMsg != "" {
		b.WriteString("\n")
		b.WriteString(errSty.Render("  └─ " + m.errorMsg))
	}

	boxSty := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#007AFF")).
		Padding(0, 1).
		Width(52)

	if m.focused {
		boxSty = boxSty.Bold(true)
	}

	return boxSty.Render(b.String())
}

func (m InputFieldModel) WithValue(value string) InputFieldModel {
	m.value = value
	m.cursorPos = utf8.RuneCountInString(value)
	m.validate()
	return m
}

func (m InputFieldModel) WithValidation(validation func(string) error) InputFieldModel {
	m.validation = validation
	m.validate()
	return m
}

func (m InputFieldModel) Focused() InputFieldModel {
	m.focused = true
	return m
}

func (m InputFieldModel) Blur() InputFieldModel {
	m.focused = false
	return m
}
