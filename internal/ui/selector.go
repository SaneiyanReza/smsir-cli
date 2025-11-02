package ui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// SelectorModel represents the mode selection model
type SelectorModel struct {
	choices  []string
	cursor   int
	Selected string
	quitting bool
	width    int
	height   int
}

// Init initializes the selector model
func (m SelectorModel) Init() tea.Cmd {
	return nil
}

// Update handles messages
func (m SelectorModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			m.quitting = true
			return m, tea.Quit

		case "enter":
			m.Selected = m.choices[m.cursor]
			return m, tea.Quit

		case "down", "j":
			if m.cursor < len(m.choices)-1 {
				m.cursor++
			}

		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		}

	default:
		return m, nil
	}

	return m, nil
}

// View renders the selector
func (m SelectorModel) View() string {
	if m.quitting {
		return "Goodbye! 👋\n"
	}

	var s strings.Builder

	// Header
	header := m.renderHeader()
	s.WriteString(header)
	s.WriteString("\n\n")

	// Instructions
	instructions := m.renderInstructions()
	s.WriteString(instructions)
	s.WriteString("\n\n")

	// Choices
	choices := m.renderChoices()
	s.WriteString(choices)

	return s.String()
}

// renderHeader renders the header
func (m SelectorModel) renderHeader() string {
	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#F59E0B")). // SMS.ir yellow
		Align(lipgloss.Center)

	subtitleStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FBBF24")). // Lighter yellow
		Align(lipgloss.Center).
		Italic(true)

	title := "📱 SMS.ir CLI"
	subtitle := "A simple message can connect worlds with a single command"

	return titleStyle.Render(title) + "\n" + subtitleStyle.Render(subtitle)
}

// renderInstructions renders instructions
func (m SelectorModel) renderInstructions() string {
	instructionStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#6B7280")).
		Align(lipgloss.Center)

	instructions := []string{
		"Use ↑/↓ or j/k to navigate",
		"Press Enter to select",
		"Press q or Ctrl+C to quit",
	}

	return instructionStyle.Render(strings.Join(instructions, " • "))
}

// renderChoices renders the choice list
func (m SelectorModel) renderChoices() string {
	var s strings.Builder

	for i, choice := range m.choices {
		cursor := " "
		if m.cursor == i {
			cursor = ">"
		}

		choiceStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("#374151"))

		if m.cursor == i {
			choiceStyle = choiceStyle.
				Foreground(lipgloss.Color("#F59E0B")). // SMS.ir yellow
				Bold(true)
		}

		choiceText := fmt.Sprintf("%s %s", cursor, choice)
		s.WriteString(choiceStyle.Render(choiceText))
		s.WriteString("\n")
	}

	return s.String()
}

// NewSelectorModel creates a new selector model
func NewSelectorModel() SelectorModel {
	return SelectorModel{
		choices: []string{
			"🔧 Configure API Key & Line Number",
			"💻 Command Line Mode",
			"🎨 Interactive Dashboard",
		},
		cursor: 0,
	}
}
