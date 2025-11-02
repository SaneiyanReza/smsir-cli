package ui

import (
	"fmt"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// StartupModel represents the startup animation model
type StartupModel struct {
	width     int
	height    int
	progress  int
	quitting  bool
	completed bool
}

// Init initializes the startup model
func (m StartupModel) Init() tea.Cmd {
	return tickCmd()
}

// Update handles messages
func (m StartupModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c", "esc":
			m.quitting = true
			return m, tea.Quit
		}
		return m, nil

	case tickMsg:
		if m.progress < 100 {
			m.progress += 2
			return m, tickCmd()
		} else {
			m.completed = true
			time.Sleep(300 * time.Millisecond)
			return m, tea.Quit
		}

	default:
		return m, nil
	}
}

// View renders the startup animation
func (m StartupModel) View() string {
	if m.quitting {
		return "Goodbye! 👋\n"
	}

	var s strings.Builder

	// ASCII Art Logo
	logo := m.renderLogo()
	s.WriteString(logo)
	s.WriteString("\n\n")

	// Tagline
	tagline := m.renderTagline()
	s.WriteString(tagline)
	s.WriteString("\n\n")

	// Progress bar
	progress := m.renderProgress()
	s.WriteString(progress)
	s.WriteString("\n\n")

	// Loading text
	loading := m.renderLoading()
	s.WriteString(loading)

	return s.String()
}

// renderLogo renders the ASCII art logo
func (m StartupModel) renderLogo() string {
	logoStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#F59E0B")). // SMS.ir yellow
		Align(lipgloss.Center)

	logo := `
    ╔══════════════════════════════════╗
    ║        📱 SMS.ir CLI 📱          ║
    ║                                  ║
    ║    A simple message can ... 💬   ║
    ╚══════════════════════════════════╝
`

	return logoStyle.Render(logo)
}

// renderTagline renders the tagline
func (m StartupModel) renderTagline() string {
	taglineStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FBBF24")). // SMS.ir light yellow
		Align(lipgloss.Center).
		Italic(true)

	tagline := "A simple message can connect worlds with a single command"

	return taglineStyle.Render(tagline)
}

// renderProgress renders the progress bar
func (m StartupModel) renderProgress() string {
	barWidth := 40
	filled := int(float64(barWidth) * float64(m.progress) / 100.0)

	barStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FFD93D")).
		Align(lipgloss.Center)

	bar := "[" + strings.Repeat("█", filled) + strings.Repeat("░", barWidth-filled) + "]"
	percentage := fmt.Sprintf(" %d%%", m.progress)

	return barStyle.Render(bar + percentage)
}

// renderLoading renders loading text
func (m StartupModel) renderLoading() string {
	loadingStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#95E1D3")).
		Align(lipgloss.Center)

	loadingTexts := []string{
		"Loading...",
		"Connecting to SMS.ir...",
		"Preparing user interface...",
		"Almost ready!",
	}

	textIndex := (m.progress / 25) % len(loadingTexts)
	if m.completed {
		textIndex = len(loadingTexts) - 1
	}

	return loadingStyle.Render(loadingTexts[textIndex])
}

// Messages
type tickMsg time.Time

// tickCmd returns a command that sends a tick message
func tickCmd() tea.Cmd {
	return tea.Tick(50*time.Millisecond, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

// NewStartupModel creates a new startup model
func NewStartupModel() StartupModel {
	return StartupModel{
		progress: 0,
	}
}
