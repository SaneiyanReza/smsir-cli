package commands

import (
	"fmt"

	"github.com/SaneiyanReza/smsir-cli/internal/api"
	"github.com/spf13/cobra"
)

// linesCmd represents the lines command
var linesCmd = &cobra.Command{
	Use:   "lines",
	Short: "Show available lines",
	Long:  `Show list of available lines for sending SMS`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client := api.NewClient(cfg)

		resp, err := client.GetLines()
		if err != nil {
			return fmt.Errorf("error getting lines: %w", err)
		}

		if !resp.IsSuccess() {
			return fmt.Errorf("API error: %s", resp.GetStatusMessage())
		}

		lines := []int64(resp.Data)

		if len(lines) == 0 {
			fmt.Println("📞 No lines found")
			return nil
		}

		fmt.Println("📞 Available Lines:")
		for i, line := range lines {
			fmt.Printf("  %d. %d\n", i+1, line)
		}

		return nil
	},
}
