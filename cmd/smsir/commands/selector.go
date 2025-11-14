package commands

import (
	"github.com/SaneiyanReza/smsir-cli/internal/ui"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
)

// selectorCmd represents the menu command
var selectorCmd = &cobra.Command{
	Use:   "menu",
	Short: "Launch interactive menu",
	Long:  `Launch interactive menu with beautiful user interface to choose between different modes: dashboard with real-time updates, configuration setup, and command line operations. This provides an easy-to-use graphical interface for navigating all SMS.ir CLI features.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		launcher := ui.NewLauncherModel()
		p := tea.NewProgram(launcher, tea.WithAltScreen())

		finalModel, err := p.Run()
		if err != nil {
			return err
		}

		// Check if help should be run after UI exits
		// finalModel should be the same LauncherModel with updated state
		if finalLauncher, ok := finalModel.(ui.LauncherModel); ok {
			if finalLauncher.ShouldRunHelp() {
				RootCmd.SetArgs([]string{"--help"})
				return RootCmd.Execute()
			}
		}

		return nil
	},
}
