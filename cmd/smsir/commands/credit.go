package commands

import (
	"fmt"

	"github.com/SaneiyanReza/smsir-cli/internal/api"
	"github.com/spf13/cobra"
)

// creditCmd represents the credit command
var creditCmd = &cobra.Command{
	Use:   "credit",
	Short: "Show current credit balance",
	Long:  `Show current credit balance of SMS.ir account`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client := api.NewClient(cfg)

		resp, err := client.GetCredit()
		if err != nil {
			return fmt.Errorf("error getting credit: %w", err)
		}

		if !resp.IsSuccess() {
			return fmt.Errorf("API error: %s", resp.GetStatusMessage())
		}

		credit := float64(resp.Data)
		fmt.Printf("💰 Current Credit: %.2f SMS\n", credit)
		return nil
	},
}
