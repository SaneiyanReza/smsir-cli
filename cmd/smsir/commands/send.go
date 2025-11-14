package commands

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/SaneiyanReza/smsir-cli/internal/api"
	"github.com/spf13/cobra"
)

var sendCmd = &cobra.Command{
	Use:   "send",
	Short: "Send SMS message",
	Long:  `Send SMS message to one or more mobile numbers`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client := api.NewClient(cfg)

		message, err := cmd.Flags().GetString("message")
		if err != nil {
			return fmt.Errorf("error getting message flag: %w", err)
		}
		if message == "" {
			return fmt.Errorf("message is required")
		}

		mobilesStr, err := cmd.Flags().GetString("to")
		if err != nil {
			return fmt.Errorf("error getting to flag: %w", err)
		}
		if mobilesStr == "" {
			return fmt.Errorf("to mobiles is required")
		}

		mobiles := strings.Split(mobilesStr, ",")
		for i := range mobiles {
			mobiles[i] = strings.TrimSpace(mobiles[i])
		}

		lineNumberStr, err := cmd.Flags().GetString("line")
		if err != nil {
			return fmt.Errorf("error getting line flag: %w", err)
		}

		var lineNumber int64
		if lineNumberStr != "" {
			lineNumber, err = strconv.ParseInt(lineNumberStr, 10, 64)
			if err != nil {
				return fmt.Errorf("invalid line number: %w", err)
			}
		} else {
			if cfg.LineNumber == "" {
				return fmt.Errorf("line number is required (use --line flag or configure it)")
			}
			lineNumber, err = strconv.ParseInt(cfg.LineNumber, 10, 64)
			if err != nil {
				return fmt.Errorf("invalid line number in config: %w", err)
			}
		}

		req := api.BulkSendRequest{
			LineNumber:  lineNumber,
			MessageText: message,
			Mobiles:     mobiles,
		}

		resp, err := client.SendBulk(req)
		if err != nil {
			return fmt.Errorf("error sending SMS: %w", err)
		}

		if !resp.IsSuccess() {
			return fmt.Errorf("API error: %s", resp.GetStatusMessage())
		}

		fmt.Printf("✅ SMS sent successfully!\n")
		fmt.Printf("📦 Pack ID: %s\n", resp.Data.PackID)
		fmt.Printf("💰 Cost: %.2f SMS\n", resp.Data.Cost)
		fmt.Printf("📱 Message IDs: %v\n", resp.Data.MessageIds)
		fmt.Printf("📊 Total messages: %d\n", len(resp.Data.MessageIds))

		return nil
	},
}

func init() {
	sendCmd.Flags().StringP("message", "m", "", "Message text to send")
	sendCmd.Flags().StringP("to", "t", "", "Comma-separated list of mobile numbers")
	sendCmd.Flags().StringP("line", "l", "", "Line number (optional, uses config if not provided)")

	sendCmd.MarkFlagRequired("message")
	sendCmd.MarkFlagRequired("to")
}
