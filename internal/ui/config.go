package ui

import (
	"fmt"
	"strings"

	"github.com/SaneiyanReza/smsir-cli/internal/config"
	"github.com/atotto/clipboard"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// ConfigModel represents the configuration model
type ConfigModel struct {
	apiKey     string
	lineNumber string
	cursor     int
	quitting   bool
	completed  bool
	width      int
	height     int
	step       int // 0: api key, 1: line number, 2: confirm
}

// Init initializes the config model
func (m ConfigModel) Init() tea.Cmd {
	return nil
}

// Update handles messages
func (m ConfigModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case tea.KeyMsg:
		// Handle paste first
		if msg.Type == tea.KeyCtrlV {
			clipboardText, err := clipboard.ReadAll()
			if err == nil && clipboardText != "" {
				// Clean the clipboard text
				cleanText := strings.TrimSpace(strings.ReplaceAll(clipboardText, "\n", ""))
				cleanText = strings.ReplaceAll(cleanText, "\r", "")

				if m.step == 0 {
					m.apiKey = cleanText
				} else if m.step == 1 {
					m.lineNumber = cleanText
				}
			}
			return m, nil
		}

		switch msg.String() {
		case "ctrl+c", "q", "esc":
			m.quitting = true
			return m, tea.Quit

		case "ctrl+v":
			// Alternative paste method
			clipboardText, err := clipboard.ReadAll()
			if err == nil && clipboardText != "" {
				cleanText := strings.TrimSpace(strings.ReplaceAll(clipboardText, "\n", ""))
				cleanText = strings.ReplaceAll(cleanText, "\r", "")

				if m.step == 0 {
					m.apiKey = cleanText
				} else if m.step == 1 {
					m.lineNumber = cleanText
				}
			}
			return m, nil

		case "enter":
			if m.step == 2 {
				// Save configuration
				cfg := &config.Config{
					APIKey:     m.apiKey,
					LineNumber: m.lineNumber,
					BaseURL:    "https://api.sms.ir/v1",
				}
				if err := cfg.SaveConfig(); err != nil {
					// Error will be handled by launcher
					m.completed = false
					m.quitting = true
					return m, tea.Quit
				}
				m.completed = true
				return m, tea.Quit
			}

			// Only advance if we have valid input
			if m.step == 0 && m.apiKey != "" {
				m.step++
			} else if m.step == 1 && m.lineNumber != "" {
				m.step++
			}
			return m, nil

		case "backspace":
			if m.step == 0 && len(m.apiKey) > 0 {
				runes := []rune(m.apiKey)
				if len(runes) > 0 {
					m.apiKey = string(runes[:len(runes)-1])
				}
			} else if m.step == 1 && len(m.lineNumber) > 0 {
				runes := []rune(m.lineNumber)
				if len(runes) > 0 {
					m.lineNumber = string(runes[:len(runes)-1])
				}
			}
			return m, nil

		default:
			// Accept printable runes, including multi-rune inserts (e.g., pastes on some terminals)
			if msg.Type == tea.KeyRunes {
				text := msg.String()
				if text != "" {
					if m.step == 0 {
						m.apiKey += text
					} else if m.step == 1 {
						m.lineNumber += text
					}
				}
				return m, nil
			}
			return m, nil
		}

	default:
		return m, nil
	}

}

// View renders the config interface
func (m ConfigModel) View() string {
	if m.quitting && !m.completed {
		return "Configuration cancelled.\n"
	}

	if m.completed {
		return ""
	}

	var s strings.Builder

	// Header
	header := m.renderHeader()
	s.WriteString(header)
	s.WriteString("\n\n")

	// Progress
	progress := m.renderProgress()
	s.WriteString(progress)
	s.WriteString("\n\n")

	// Content
	content := m.renderContent()
	s.WriteString(content)
	s.WriteString("\n\n")

	// Instructions
	instructions := m.renderInstructions()
	s.WriteString(instructions)

	return s.String()
}

// renderHeader renders the header
func (m ConfigModel) renderHeader() string {
	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#f7bd60")).
		Align(lipgloss.Center)

	title := "🔧 Configuration Setup"
	return titleStyle.Render(title)
}

// renderProgress renders progress indicator
func (m ConfigModel) renderProgress() string {
	steps := []string{"API Key", "Line Number", "Confirm"}

	var progress []string
	for i, step := range steps {
		if i <= m.step {
			progress = append(progress, fmt.Sprintf("✓ %s", step))
		} else {
			progress = append(progress, fmt.Sprintf("○ %s", step))
		}
	}

	progressStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#9CA3AF")).
		Align(lipgloss.Center)

	return progressStyle.Render(strings.Join(progress, " → "))
}

// renderContent renders the main content
func (m ConfigModel) renderContent() string {
	boxStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#f7bd60")).
		Padding(1, 2).
		Width(m.width - 4)

	var content string

	switch m.step {
	case 0:
		content = m.renderAPIKeyStep()
	case 1:
		content = m.renderLineNumberStep()
	case 2:
		content = m.renderConfirmStep()
	}

	return boxStyle.Render(content)
}

// renderAPIKeyStep renders the API key input step
func (m ConfigModel) renderAPIKeyStep() string {
	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#f7bd60"))

	inputStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#ffffff")).
		Bold(true)

	placeholderStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#9CA3AF")).
		Italic(true)

	title := "Enter your SMS.ir API Key:"
	var input string
	if m.apiKey == "" {
		input = placeholderStyle.Render("Type here or press Ctrl+V to paste...")
	} else {
		input = inputStyle.Render(m.apiKey)
	}

	return titleStyle.Render(title) + "\n\n" + input
}

// renderLineNumberStep renders the line number input step
func (m ConfigModel) renderLineNumberStep() string {
	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#f7bd60"))

	inputStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#ffffff")).
		Bold(true)

	placeholderStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#9CA3AF")).
		Italic(true)

	title := "Enter your Line Number:"
	var input string
	if m.lineNumber == "" {
		input = placeholderStyle.Render("Type here or press Ctrl+V to paste...")
	} else {
		input = inputStyle.Render(m.lineNumber)
	}

	return titleStyle.Render(title) + "\n\n" + input
}

// renderConfirmStep renders the confirmation step
func (m ConfigModel) renderConfirmStep() string {
	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#f7bd60"))

	infoStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#ffffff"))

	title := "Confirm Configuration:"
	apiKey := fmt.Sprintf("API Key: %s", maskString(m.apiKey))
	lineNumber := fmt.Sprintf("Line Number: %s", m.lineNumber)

	return titleStyle.Render(title) + "\n\n" +
		infoStyle.Render(apiKey) + "\n" +
		infoStyle.Render(lineNumber)
}

// renderInstructions renders instructions
func (m ConfigModel) renderInstructions() string {
	instructionStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#9CA3AF")).
		Align(lipgloss.Center)

	var instructions []string
	if m.step < 2 {
		instructions = []string{
			"Type your information and press Enter to continue",
			"Press Ctrl+V to paste from clipboard",
			"Press q or Ctrl+C to cancel",
		}
	} else {
		instructions = []string{
			"Press Enter to save configuration",
			"Press q or Ctrl+C to cancel",
		}
	}

	return instructionStyle.Render(strings.Join(instructions, " • "))
}

// maskString masks sensitive information
func maskString(s string) string {
	if len(s) <= 8 {
		return "****"
	}
	return s[:4] + "****" + s[len(s)-4:]
}

// NewConfigModel creates a new config model
func NewConfigModel() ConfigModel {
	return ConfigModel{
		cursor: 0,
		step:   0,
	}
}
