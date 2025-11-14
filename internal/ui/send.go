package ui

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/SaneiyanReza/smsir-cli/internal/api"
	"github.com/SaneiyanReza/smsir-cli/internal/config"
	"github.com/atotto/clipboard"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// SendModel represents the send SMS model
type SendModel struct {
	client      *api.Client
	config      *config.Config
	messageText string
	mobiles     string
	lineNumber  string
	quitting    bool
	completed   bool
	success     bool
	result      *api.BulkSendResponse
	err         error
	width       int
	height      int
	step        int // 0: message, 1: mobiles, 2: line number (optional), 3: confirm
}

// Init initializes the send model
func (m SendModel) Init() tea.Cmd {
	return nil
}

// Update handles messages
func (m SendModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
				cleanText := strings.TrimSpace(strings.ReplaceAll(clipboardText, "\n", ""))
				cleanText = strings.ReplaceAll(cleanText, "\r", "")

				if m.step == 0 {
					m.messageText = cleanText
				} else if m.step == 1 {
					m.mobiles = cleanText
				} else if m.step == 2 {
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
			clipboardText, err := clipboard.ReadAll()
			if err == nil && clipboardText != "" {
				cleanText := strings.TrimSpace(strings.ReplaceAll(clipboardText, "\n", ""))
				cleanText = strings.ReplaceAll(cleanText, "\r", "")

				if m.step == 0 {
					m.messageText = cleanText
				} else if m.step == 1 {
					m.mobiles = cleanText
				} else if m.step == 2 {
					m.lineNumber = cleanText
				}
			}
			return m, nil

		case "enter":
			if m.step == 3 {
				// Send SMS
				return m, m.sendSMS()
			}

			// Skip line number step if empty and config has line number
			if m.step == 2 {
				if m.lineNumber == "" && m.config.LineNumber != "" {
					m.step++
					return m, nil
				}
			}

			// Only advance if we have valid input
			if m.step == 0 && m.messageText != "" {
				m.step++
			} else if m.step == 1 && m.mobiles != "" {
				m.step++
			} else if m.step == 2 {
				m.step++
			}
			return m, nil

		case "backspace":
			if m.step == 0 && len(m.messageText) > 0 {
				runes := []rune(m.messageText)
				if len(runes) > 0 {
					m.messageText = string(runes[:len(runes)-1])
				}
			} else if m.step == 1 && len(m.mobiles) > 0 {
				runes := []rune(m.mobiles)
				if len(runes) > 0 {
					m.mobiles = string(runes[:len(runes)-1])
				}
			} else if m.step == 2 && len(m.lineNumber) > 0 {
				runes := []rune(m.lineNumber)
				if len(runes) > 0 {
					m.lineNumber = string(runes[:len(runes)-1])
				}
			}
			return m, nil

		default:
			if msg.Type == tea.KeyRunes {
				text := msg.String()
				if text != "" {
					if m.step == 0 {
						m.messageText += text
					} else if m.step == 1 {
						m.mobiles += text
					} else if m.step == 2 {
						m.lineNumber += text
					}
				}
				return m, nil
			}
			return m, nil
		}

	case sendSuccessMsg:
		m.success = true
		m.result = msg.result
		m.completed = true
		return m, nil

	case sendErrorMsg:
		m.err = msg.err
		m.completed = true
		return m, nil

	default:
		return m, nil
	}
}

// View renders the send interface
func (m SendModel) View() string {
	if m.quitting && !m.completed {
		return "SMS sending cancelled.\n"
	}

	if m.completed {
		if m.success {
			return m.renderSuccess()
		}
		return m.renderError()
	}

	var s strings.Builder

	header := m.renderHeader()
	s.WriteString(header)
	s.WriteString("\n\n")

	progress := m.renderProgress()
	s.WriteString(progress)
	s.WriteString("\n\n")

	content := m.renderContent()
	s.WriteString(content)
	s.WriteString("\n\n")

	instructions := m.renderInstructions()
	s.WriteString(instructions)

	return s.String()
}

// renderHeader renders the header
func (m SendModel) renderHeader() string {
	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#f7bd60")).
		Align(lipgloss.Center)

	title := "📤 Send SMS"
	return titleStyle.Render(title)
}

// renderProgress renders progress indicator
func (m SendModel) renderProgress() string {
	steps := []string{"Message", "Mobiles", "Line Number", "Confirm"}

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
func (m SendModel) renderContent() string {
	boxStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#f7bd60")).
		Padding(1, 2).
		Width(m.width - 4)

	var content string

	switch m.step {
	case 0:
		content = m.renderMessageStep()
	case 1:
		content = m.renderMobilesStep()
	case 2:
		content = m.renderLineNumberStep()
	case 3:
		content = m.renderConfirmStep()
	}

	return boxStyle.Render(content)
}

// renderMessageStep renders the message input step
func (m SendModel) renderMessageStep() string {
	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#f7bd60"))

	inputStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#ffffff")).
		Bold(true)

	placeholderStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#9CA3AF")).
		Italic(true)

	title := "Enter your message text:"
	var input string
	if m.messageText == "" {
		input = placeholderStyle.Render("Type here or press Ctrl+V to paste...")
	} else {
		input = inputStyle.Render(m.messageText)
	}

	return titleStyle.Render(title) + "\n\n" + input
}

// renderMobilesStep renders the mobiles input step
func (m SendModel) renderMobilesStep() string {
	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#f7bd60"))

	inputStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#ffffff")).
		Bold(true)

	placeholderStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#9CA3AF")).
		Italic(true)

	title := "Enter mobile numbers (comma-separated):"
	var input string
	if m.mobiles == "" {
		input = placeholderStyle.Render("e.g., 09120000000,09121111111")
	} else {
		input = inputStyle.Render(m.mobiles)
	}

	return titleStyle.Render(title) + "\n\n" + input
}

// renderLineNumberStep renders the line number input step
func (m SendModel) renderLineNumberStep() string {
	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#f7bd60"))

	inputStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#ffffff")).
		Bold(true)

	placeholderStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#9CA3AF")).
		Italic(true)

	defaultLine := ""
	if m.config.LineNumber != "" {
		defaultLine = fmt.Sprintf(" (Press Enter to use: %s)", m.config.LineNumber)
	}

	title := fmt.Sprintf("Enter line number%s:", defaultLine)
	var input string
	if m.lineNumber == "" {
		input = placeholderStyle.Render("Leave empty to use configured line number")
	} else {
		input = inputStyle.Render(m.lineNumber)
	}

	return titleStyle.Render(title) + "\n\n" + input
}

// renderConfirmStep renders the confirmation step
func (m SendModel) renderConfirmStep() string {
	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#f7bd60"))

	infoStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#ffffff"))

	lineNumber := m.lineNumber
	if lineNumber == "" {
		lineNumber = m.config.LineNumber
		if lineNumber == "" {
			lineNumber = "Not set"
		}
	}

	mobilesList := strings.Split(m.mobiles, ",")
	for i := range mobilesList {
		mobilesList[i] = strings.TrimSpace(mobilesList[i])
	}

	title := "Confirm and Send:"
	message := fmt.Sprintf("Message: %s", m.messageText)
	mobiles := fmt.Sprintf("Mobiles: %s", strings.Join(mobilesList, ", "))
	line := fmt.Sprintf("Line Number: %s", lineNumber)

	return titleStyle.Render(title) + "\n\n" +
		infoStyle.Render(message) + "\n" +
		infoStyle.Render(mobiles) + "\n" +
		infoStyle.Render(line)
}

// renderInstructions renders instructions
func (m SendModel) renderInstructions() string {
	instructionStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#9CA3AF")).
		Align(lipgloss.Center)

	var instructions []string
	if m.step < 3 {
		instructions = []string{
			"Type your information and press Enter to continue",
			"Press Ctrl+V to paste from clipboard",
			"Press q or Ctrl+C to cancel",
		}
	} else {
		instructions = []string{
			"Press Enter to send SMS",
			"Press q or Ctrl+C to cancel",
		}
	}

	return instructionStyle.Render(strings.Join(instructions, " • "))
}

// renderSuccess renders success message
func (m SendModel) renderSuccess() string {
	boxStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#f7bd60")).
		Padding(1, 2).
		Width(m.width - 4)

	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#f7bd60"))

	infoStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#ffffff"))

	content := titleStyle.Render("✅ SMS sent successfully!") + "\n\n" +
		infoStyle.Render(fmt.Sprintf("📦 Pack ID: %s", m.result.PackID)) + "\n" +
		infoStyle.Render(fmt.Sprintf("💰 Cost: %.2f SMS", m.result.Cost)) + "\n" +
		infoStyle.Render(fmt.Sprintf("📱 Message IDs: %v", m.result.MessageIds)) + "\n" +
		infoStyle.Render(fmt.Sprintf("📊 Total messages: %d", len(m.result.MessageIds))) + "\n\n" +
		infoStyle.Render("Press q or Ctrl+C to exit...")

	return boxStyle.Render(content)
}

// renderError renders error message
func (m SendModel) renderError() string {
	boxStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#FF6B6B")).
		Padding(1, 2).
		Width(m.width - 4)

	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FF6B6B"))

	errorStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#ffffff"))

	content := titleStyle.Render("❌ Error sending SMS") + "\n\n" +
		errorStyle.Render(fmt.Sprintf("Error: %v", m.err)) + "\n\n" +
		errorStyle.Render("Press q or Ctrl+C to exit...")

	return boxStyle.Render(content)
}

// Messages
type sendSuccessMsg struct {
	result *api.BulkSendResponse
}

type sendErrorMsg struct {
	err error
}

// sendSMS sends the SMS
func (m SendModel) sendSMS() tea.Cmd {
	return func() tea.Msg {
		mobilesList := strings.Split(m.mobiles, ",")
		for i := range mobilesList {
			mobilesList[i] = strings.TrimSpace(mobilesList[i])
		}

		lineNumberStr := m.lineNumber
		if lineNumberStr == "" {
			lineNumberStr = m.config.LineNumber
		}

		if lineNumberStr == "" {
			return sendErrorMsg{err: fmt.Errorf("line number is required")}
		}

		lineNumber, err := strconv.ParseInt(lineNumberStr, 10, 64)
		if err != nil {
			return sendErrorMsg{err: fmt.Errorf("invalid line number: %w", err)}
		}

		req := api.BulkSendRequest{
			LineNumber:  lineNumber,
			MessageText: m.messageText,
			Mobiles:     mobilesList,
		}

		resp, err := m.client.SendBulk(req)
		if err != nil {
			return sendErrorMsg{err: err}
		}

		if !resp.IsSuccess() {
			return sendErrorMsg{err: fmt.Errorf("API error: %s", resp.GetStatusMessage())}
		}

		return sendSuccessMsg{result: &resp.Data}
	}
}

// NewSendModel creates a new send model
func NewSendModel(client *api.Client, cfg *config.Config) SendModel {
	return SendModel{
		client: client,
		config: cfg,
		step:   0,
	}
}
