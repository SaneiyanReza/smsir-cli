package commands

import (
	"fmt"

	"github.com/SaneiyanReza/smsir-cli/internal/config"
	"github.com/spf13/cobra"
)

// configCmd represents the config command
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Configuration management",
	Long:  `Manage API Key and line number configuration`,
}

// configSetCmd represents the config set command
var configSetCmd = &cobra.Command{
	Use:   "set",
	Short: "Set API Key and line number",
	Long:  `Set API Key and line number for SMS.ir connection`,
	RunE: func(cmd *cobra.Command, args []string) error {
		apiKey, _ := cmd.Flags().GetString("api-key")
		lineNumber, _ := cmd.Flags().GetString("line")

		if apiKey == "" {
			return fmt.Errorf("api key is required")
		}
		if lineNumber == "" {
			return fmt.Errorf("line number is required")
		}

		// Load or create config
		var err error
		cfg, err = config.LoadConfig()
		if err != nil {
			return fmt.Errorf("error loading configuration: %w", err)
		}

		cfg.APIKey = apiKey
		cfg.LineNumber = lineNumber

		if err := cfg.SaveConfig(); err != nil {
			return fmt.Errorf("error saving configuration: %w", err)
		}

		fmt.Println("✅ Configuration saved successfully")
		return nil
	},
}

// configShowCmd represents the config show command
var configShowCmd = &cobra.Command{
	Use:   "show",
	Short: "Show current configuration",
	Long:  `Show current configuration including API Key and line number`,
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Printf("API Key: %s\n", maskString(cfg.APIKey))
		fmt.Printf("Line Number: %s\n", cfg.LineNumber)
		fmt.Printf("Base URL: %s\n", cfg.BaseURL)
		return nil
	},
}

// configValidateCmd represents the config validate command
var configValidateCmd = &cobra.Command{
	Use:   "validate",
	Short: "Validate configuration",
	Long:  `Validate current configuration`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := cfg.Validate(); err != nil {
			return fmt.Errorf("invalid configuration: %w", err)
		}

		fmt.Println("✅ Configuration is valid")
		return nil
	},
}

func init() {
	// Note: configCmd is added to rootCmd in root.go setupCommands() to control order
	// Only add sub-commands here
	configCmd.AddCommand(configSetCmd)
	configCmd.AddCommand(configShowCmd)
	configCmd.AddCommand(configValidateCmd)

	configSetCmd.Flags().String("api-key", "", "API Key from SMS.ir panel")
	configSetCmd.Flags().String("line", "", "Line number")
	configSetCmd.MarkFlagRequired("api-key")
	configSetCmd.MarkFlagRequired("line")
}

// maskString masks sensitive information
func maskString(s string) string {
	if len(s) <= 8 {
		return "****"
	}
	return s[:4] + "****" + s[len(s)-4:]
}
