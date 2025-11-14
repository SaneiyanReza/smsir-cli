package ui

import (
	"os"
	"strings"

	"github.com/SaneiyanReza/smsir-cli/internal/api"
	"github.com/SaneiyanReza/smsir-cli/internal/config"
	tea "github.com/charmbracelet/bubbletea"
)

// Launcher state constants
const (
	stateStartup   = "startup"
	stateSelector  = "selector"
	stateConfig    = "config"
	stateDashboard = "dashboard"
	stateSend      = "send"
	stateDone      = "done"
	stateHelp      = "help"
	stateExit      = "exit"
)

// LauncherModel represents the main launcher that runs startup then selector
type LauncherModel struct {
	state         string // Use state constants above
	startup       StartupModel
	selector      SelectorModel
	config        ConfigModel
	dashboard     Model
	send          SendModel
	width         int
	height        int
	helpOutput    string
	shouldRunHelp bool // Flag to indicate help should be run after UI exits
}

// Init initializes the launcher model
func (m LauncherModel) Init() tea.Cmd {
	return m.startup.Init()
}

// Update handles messages
func (m LauncherModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.startup.width = msg.Width
		m.startup.height = msg.Height
		m.selector.width = msg.Width
		m.selector.height = msg.Height
		if m.config.width > 0 || m.state == stateConfig {
			m.config.width = msg.Width
			m.config.height = msg.Height
		}
		if m.dashboard.width > 0 || m.state == stateDashboard {
			m.dashboard.width = msg.Width
			m.dashboard.height = msg.Height
		}
		if m.send.width > 0 || m.state == stateSend {
			m.send.width = msg.Width
			m.send.height = msg.Height
		}
		return m, nil

	default:
		switch m.state {
		case stateStartup:
			startupModel, cmd := m.startup.Update(msg)
			if sm, ok := startupModel.(StartupModel); ok {
				m.startup = sm
			}

			// If startup is completed, move to selector
			if m.startup.completed {
				m.state = stateSelector
				return m, m.selector.Init()
			}

			return m, cmd

		case stateSelector:
			selectorModel, cmd := m.selector.Update(msg)
			if sm, ok := selectorModel.(SelectorModel); ok {
				m.selector = sm
			}

			// If selector is done, handle the selection
			if m.selector.Selected != "" {
				return m.handleSelection()
			}

			if m.selector.quitting {
				m.state = stateDone
				return m, tea.Quit
			}

			return m, cmd

		case stateConfig:
			configModel, cmd := m.config.Update(msg)
			if cm, ok := configModel.(ConfigModel); ok {
				m.config = cm
			}

			// If config is completed, show success message briefly then return to selector
			if m.config.completed {
				// Success message will be cleared when we return to selector
				m.state = stateSelector
				m.selector = NewSelectorModel()
				m.selector.width = m.width
				m.selector.height = m.height
				return m, m.selector.Init()
			}

			if m.config.quitting {
				m.state = stateSelector
				m.selector = NewSelectorModel()
				m.selector.width = m.width
				m.selector.height = m.height
				return m, m.selector.Init()
			}

			return m, cmd

		case stateDashboard:
			dashboardModel, cmd := m.dashboard.Update(msg)
			if dm, ok := dashboardModel.(Model); ok {
				m.dashboard = dm
			}

			// If dashboard is quitting, return to selector
			if m.dashboard.quitting {
				m.state = stateSelector
				m.selector = NewSelectorModel()
				m.selector.width = m.width
				m.selector.height = m.height
				return m, m.selector.Init()
			}

			return m, cmd

		case stateSend:
			sendModel, cmd := m.send.Update(msg)
			if sm, ok := sendModel.(SendModel); ok {
				m.send = sm
			}

			if m.send.quitting {
				m.state = stateSelector
				m.selector = NewSelectorModel()
				m.selector.width = m.width
				m.selector.height = m.height
				return m, m.selector.Init()
			}

			return m, cmd

		case stateHelp:
			// Wait for any key press to exit
			if keyMsg, ok := msg.(tea.KeyMsg); ok {
				switch keyMsg.String() {
				case "q", "ctrl+c", "enter", "esc":
					return m, tea.Quit
				}
			}
			return m, nil

		case stateDone, stateExit:
			return m, tea.Quit
		}
	}

	return m, nil
}

// View renders the launcher
func (m LauncherModel) View() string {
	switch m.state {
	case stateStartup:
		return m.startup.View()
	case stateSelector:
		return m.selector.View()
	case stateConfig:
		return m.config.View()
	case stateDashboard:
		return m.dashboard.View()
	case stateSend:
		return m.send.View()
	case stateHelp:
		return m.helpOutput
	case stateDone, stateExit:
		return ""
	default:
		return ""
	}
}

// handleSelection handles the user's selection and transitions to appropriate state
func (m *LauncherModel) handleSelection() (tea.Model, tea.Cmd) {
	switch m.selector.Selected {
	case "🔧 Configure API Key & Line Number":
		// Transition to config state
		m.state = stateConfig
		m.config = NewConfigModel()
		m.config.width = m.width
		m.config.height = m.height
		return m, nil

	case "📤 Send SMS":
		// Load config and transition to send state
		cfg, err := config.LoadConfig()
		if err != nil {
			m.state = stateSelector
			m.selector = NewSelectorModel()
			m.selector.width = m.width
			m.selector.height = m.height
			return m, m.selector.Init()
		}

		if err := cfg.Validate(); err != nil {
			m.state = stateSelector
			m.selector = NewSelectorModel()
			m.selector.width = m.width
			m.selector.height = m.height
			return m, m.selector.Init()
		}

		client := api.NewClient(cfg)
		m.state = stateSend
		m.send = NewSendModel(client, cfg)
		m.send.width = m.width
		m.send.height = m.height
		return m, m.send.Init()

	case "💻 Command Line Mode":
		// Exit UI and run help command
		m.shouldRunHelp = true
		m.state = stateExit
		return m, tea.Quit

	case "📊 Dashboard":
		// Load config and transition to dashboard state
		cfg, err := config.LoadConfig()
		if err != nil {
			m.state = stateSelector
			m.selector = NewSelectorModel()
			m.selector.width = m.width
			m.selector.height = m.height
			// Show error message (will be displayed when selector is shown)
			return m, m.selector.Init()
		}

		if err := cfg.Validate(); err != nil {
			m.state = stateSelector
			m.selector = NewSelectorModel()
			m.selector.width = m.width
			m.selector.height = m.height
			return m, m.selector.Init()
		}

		client := api.NewClient(cfg)
		m.state = stateDashboard
		m.dashboard = NewModel(client, cfg)
		m.dashboard.width = m.width
		m.dashboard.height = m.height
		return m, m.dashboard.Init()
	}

	// Default: return to selector
	m.state = stateSelector
	m.selector = NewSelectorModel()
	m.selector.width = m.width
	m.selector.height = m.height
	return m, m.selector.Init()
}

// getHelpOutput gets the help output by running the command
func (m *LauncherModel) getHelpOutput() string {
	// Get help output from the actual command
	cmd := os.Args[0]
	if len(os.Args) > 0 {
		// Try to get the actual smsir command path
		cmd = "smsir"
	}

	var output strings.Builder
	output.WriteString("SMS.ir CLI - Command Line Mode\n\n")
	output.WriteString("Available Commands:\n")
	output.WriteString("  config    Configuration management\n")
	output.WriteString("  send      Send SMS message\n")
	output.WriteString("  credit    Show current credit balance\n")
	output.WriteString("  lines     Show available lines\n")
	output.WriteString("  menu      Launch interactive menu\n")
	output.WriteString("  help      Help about any command\n\n")
	output.WriteString("Usage:\n")
	output.WriteString("  " + cmd + " [command]\n\n")
	output.WriteString("Examples:\n")
	output.WriteString("  " + cmd + " config set --api-key YOUR_KEY --line YOUR_LINE\n")
	output.WriteString("  " + cmd + " send -m \"Hello\" -t \"09120000000,09121111111\"\n\n")
	output.WriteString("  " + cmd + " credit\n")
	output.WriteString("  " + cmd + " lines\n")
	output.WriteString("For more information, run: " + cmd + " --help\n\n")
	output.WriteString("Press any key to exit...")

	return output.String()
}

// ShouldRunHelp returns whether help should be run after UI exits
func (m LauncherModel) ShouldRunHelp() bool {
	return m.shouldRunHelp
}

// NewLauncherModel creates a new launcher model
func NewLauncherModel() LauncherModel {
	return LauncherModel{
		state:    stateStartup,
		startup:  NewStartupModel(),
		selector: NewSelectorModel(),
	}
}
