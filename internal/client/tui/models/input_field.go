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
	Value        string
	Validation   func(string) error
	IsPassword   bool
}

type InputFieldFocusedMsg struct{ Index int }
type InputFieldValueChangeMsg struct{ Index int; Value string }

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
		case "right", "ctrl+l":
			if m.cursorPos < utf8.RuneCountInString(m.value) {
				m.cursorPos++
			}
		case "backspace", "ctrl+h":
			if m.cursorPos > 0 {
				runes := []rune(m.value)
				m.value = string(append(runes[:m.cursorPos-1], runes[m.cursorPos:]...))
				m.cursorPos--
				m.validate()
			}
		case "delete", "ctrl+d":
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
	baseStyle := lipgloss.NewStyle().
		Padding(0, 2).
		Width(50)

	labelStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#007AFF")).
		Bold(true).
		MarginBottom(0)

	placeholderStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#8E8E93")).
		Italic(true)

	valueStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FFFFFF"))

	errorStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FF3B30")).
		MarginTop(0)

	if m.isPassword {
		valueStyle = valueStyle.Copy().
			Foreground(lipgloss.Color("#8E8E93")) // Muted color for password
	} else {
		if m.value == "" && m.placeholder != "" {
			return lipgloss.JoinHorizontal(
				lipgloss.Left,
				baseStyle.Render(""),
				labelStyle.Render(m.label + ":"),
				placeholderStyle.Render(m.placeholder),
			)
		}
	}

	var label string
	if m.value != "" {
		label = labelStyle.Render(m.label + ":")
	}

	// Show cursor position
	valueRunes := []rune(m.value)
	if m.isPassword {
		valueRunes = []rune(strings.Repeat("•", utf8.RuneCountInString(m.value)))
	}
	cursorLine := strings.Repeat(" ", m.cursorPos)
	if m.cursorPos < len(valueRunes) {
		cursorLine += "▊"
	}

	valueDisplay := string(valueRunes[:m.cursorPos]) + cursorLine + string(valueRunes[m.cursorPos:])

	var content string
	if m.errorMsg != "" {
		content = lipgloss.JoinVertical(
			lipgloss.Left,
			baseStyle.Render(""),
			label,
			valueStyle.Render(valueDisplay),
			"",
			errorStyle.Render("  └─ " + m.errorMsg),
		)
	} else {
		content = lipgloss.JoinVertical(
			lipgloss.Left,
			baseStyle.Render(""),
			label,
			valueStyle.Render(valueDisplay),
		)
	}

	boxStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#007AFF")).
		Padding(0, 1).
		Width(52)

	if m.focused {
		boxStyle = boxStyle.Copy().
			Foreground(lipgloss.Color("#007AFF")).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#007AFF")).
			Bold(true)
	}

	return boxStyle.Render(content)
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
