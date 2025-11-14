package commands

import (
	"fmt"
	"os"

	"github.com/SaneiyanReza/smsir-cli/internal/config"
	"github.com/spf13/cobra"
)

var (
	cfg     *config.Config
	RootCmd = &cobra.Command{
		Use:   "smsir",
		Short: "SMS.ir CLI - A simple message can connect worlds with a single command",
		Long: `SMS.ir CLI is a professional command-line tool for interacting with SMS.ir APIs.

This tool provides both interactive and command-line interfaces for:
• Sending SMS messages
• Interactive dashboard with real-time updates
• Checking setting

Quick Start:
  smsir config                    # Set API credentials
  smsir send                      # Send SMS message
  smsir credit                    # Check your balance
  smsir lines                     # View available lines
  smsir menu                      # Launch interactive menu`,
	}
)

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() {
	err := RootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	RootCmd.PersistentPreRunE = func(cmd *cobra.Command, args []string) error {
		// Skip config loading for config set command
		if cmd.Name() == "set" {
			return nil
		}

		var err error
		cfg, err = config.LoadConfig()
		if err != nil {
			return fmt.Errorf("error loading configuration: %w", err)
		}
		return nil
	}

	RootCmd.PersistentFlags().BoolP("verbose", "v", false, "show more details")

	// Disable completion command
	RootCmd.CompletionOptions.DisableDefaultCmd = true

	// Disable alphabetical sorting to maintain custom order
	cobra.EnableCommandSorting = false

	// Commands are added in setupCommands() to maintain custom order
	setupCommands()
}

// setupCommands adds all commands to rootCmd in the desired order
func setupCommands() {
	// Configuration first (most important)
	RootCmd.AddCommand(configCmd)

	// Send command
	RootCmd.AddCommand(sendCmd)

	// Then credit and lines
	RootCmd.AddCommand(creditCmd)
	RootCmd.AddCommand(linesCmd)

	// Menu last (UI access point)
	RootCmd.AddCommand(selectorCmd)
}
