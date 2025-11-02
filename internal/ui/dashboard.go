package ui

import (
	"fmt"
	"strings"

	"github.com/SaneiyanReza/smsir-cli/internal/api"
	"github.com/SaneiyanReza/smsir-cli/internal/config"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Model represents the TUI model
type Model struct {
	client   *api.Client
	config   *config.Config
	credit   float64
	lines    []int64
	loading  bool
	err      error
	width    int
	height   int
	quitting bool
}

// Init initializes the model
func (m Model) Init() tea.Cmd {
	return tea.Batch(
		loadCredit(m.client),
		loadLines(m.client),
	)
}

// Update handles messages
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			m.quitting = true
			return m, tea.Quit
		case "r":
			return m, tea.Batch(
				loadCredit(m.client),
				loadLines(m.client),
			)
		}
		return m, nil

	case creditMsg:
		m.credit = float64(msg)
		m.loading = false
		return m, nil

	case linesMsg:
		m.lines = []int64(msg)
		m.loading = false
		return m, nil

	case errMsg:
		m.err = msg
		m.loading = false
		return m, nil

	default:
		return m, nil
	}
}

// View renders the UI
func (m Model) View() string {
	if m.quitting {
		return ""
	}

	var s strings.Builder

	// Header
	header := m.renderHeader()
	s.WriteString(header)
	s.WriteString("\n\n")

	// Content
	if m.loading {
		s.WriteString(m.renderLoading())
	} else if m.err != nil {
		s.WriteString(m.renderError())
	} else {
		s.WriteString(m.renderContent())
	}

	return s.String()
}

// renderHeader renders the application header
func (m Model) renderHeader() string {
	title := "📱 SMS.ir CLI Dashboard"
	subtitle := "A simple message can connect worlds with a single command"

	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#F59E0B")). // SMS.ir yellow
		Align(lipgloss.Center)

	subtitleStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FBBF24")). // Lighter yellow
		Align(lipgloss.Center).
		Italic(true)

	return titleStyle.Render(title) + "\n" + subtitleStyle.Render(subtitle)
}

// renderLoading renders loading state
func (m Model) renderLoading() string {
	loadingStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FFD93D")).
		Align(lipgloss.Center)

	return loadingStyle.Render("Loading... ⏳")
}

// renderError renders error state
func (m Model) renderError() string {
	errorStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FF6B6B")).
		Align(lipgloss.Center).
		Border(lipgloss.RoundedBorder()).
		Padding(1, 2)

	return errorStyle.Render(fmt.Sprintf("Error: %v", m.err))
}

// renderContent renders the main content
func (m Model) renderContent() string {
	var s strings.Builder

	// Credit section
	creditBox := m.renderCreditBox()
	s.WriteString(creditBox)
	s.WriteString("\n\n")

	// Lines section
	linesBox := m.renderLinesBox()
	s.WriteString(linesBox)
	s.WriteString("\n\n")

	// Instructions
	instructions := m.renderInstructions()
	s.WriteString(instructions)

	return s.String()
}

// renderCreditBox renders the credit information box
func (m Model) renderCreditBox() string {
	boxStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#F59E0B")). // SMS.ir yellow
		Padding(1, 2).
		Width(m.width - 4)

	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#F59E0B")) // SMS.ir yellow

	valueStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FFD93D")).
		Bold(true)

	title := "💰 Current Credit"
	value := fmt.Sprintf("%.2f SMS", m.credit)

	content := titleStyle.Render(title) + "\n" + valueStyle.Render(value)

	return boxStyle.Render(content)
}

// renderLinesBox renders the lines information box
func (m Model) renderLinesBox() string {
	boxStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#10B981")). // Green for lines
		Padding(1, 2).
		Width(m.width - 4)

	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#059669")) // Dark green

	lineStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#95E1D3"))

	title := "📞 Available Lines"

	var lines []string
	for _, line := range m.lines {
		lines = append(lines, fmt.Sprintf("%d", line))
	}

	content := titleStyle.Render(title) + "\n"
	if len(lines) > 0 {
		content += strings.Join(lines, "\n")
	} else {
		content += lineStyle.Render("No lines found")
	}

	return boxStyle.Render(content)
}

// renderInstructions renders usage instructions
func (m Model) renderInstructions() string {
	instructionStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#95E1D3")).
		Align(lipgloss.Center)

	instructions := []string{
		"Commands:",
		"r - Refresh data",
		"q - Quit",
	}

	return instructionStyle.Render(strings.Join(instructions, " | "))
}

// Messages
type creditMsg float64
type linesMsg []int64
type errMsg error

// loadCredit loads credit information
func loadCredit(client *api.Client) tea.Cmd {
	return func() tea.Msg {
		resp, err := client.GetCredit()
		if err != nil {
			return errMsg(err)
		}
		return creditMsg(resp.Data)
	}
}

// loadLines loads lines information
func loadLines(client *api.Client) tea.Cmd {
	return func() tea.Msg {
		resp, err := client.GetLines()
		if err != nil {
			return errMsg(err)
		}
		return linesMsg(resp.Data)
	}
}

// NewModel creates a new TUI model
func NewModel(client *api.Client, config *config.Config) Model {
	return Model{
		client:  client,
		config:  config,
		loading: true,
	}
}
