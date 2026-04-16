package cmd

import (
	"context"
	"fmt"
	"time"

	"github.com/spf13/cobra"

	"beancount.io/beancount-cli/generated"
	"beancount.io/beancount-cli/internal/utils"
)

var eventAddCmd = &cobra.Command{
	Use:   "add",
	Short: "Add an event directive to a ledger",
	Long: `Add an event directive to a Beancount ledger.

Example:
  beancount-cli event add \
    --ledger user/mybook \
    --type "location" \
    --description "New York" \
    --date 2024-01-01`,
	RunE: runEventAdd,
}

var (
	eventLedger      string
	eventDate        string
	eventType        string
	eventDescription string
)

func init() {
	eventCmd.AddCommand(eventAddCmd)
	eventAddCmd.Flags().StringVar(&eventLedger, "ledger", "", "Ledger full name (required)")
	eventAddCmd.Flags().StringVar(&eventType, "type", "", "Event type (required)")
	eventAddCmd.Flags().StringVar(&eventDescription, "description", "", "Event description (required)")
	eventAddCmd.Flags().StringVar(&eventDate, "date", "", "Event date in YYYY-MM-DD format (default: today)")
	_ = eventAddCmd.MarkFlagRequired("ledger")
	_ = eventAddCmd.MarkFlagRequired("type")
	_ = eventAddCmd.MarkFlagRequired("description")
}

func buildEventInput(date, eventType, description string) (generated.LedgerEventInput, error) {
	if eventType == "" {
		return generated.LedgerEventInput{}, fmt.Errorf("type is required")
	}
	if description == "" {
		return generated.LedgerEventInput{}, fmt.Errorf("description is required")
	}
	if date == "" {
		date = time.Now().Format("2006-01-02")
	}
	return generated.LedgerEventInput{
		Date:        date,
		Type:        eventType,
		Description: description,
	}, nil
}

func runEventAdd(cmd *cobra.Command, _ []string) error {
	event, err := buildEventInput(eventDate, eventType, eventDescription)
	if err != nil {
		return err
	}
	client, err := utils.NewAuthedClient()
	if err != nil {
		return err
	}
	resp, err := generated.AddEntryEvent(context.Background(), client, utils.LedgerID(eventLedger), event)
	if err != nil {
		return fmt.Errorf("failed to add event directive: %w", err)
	}
	result := resp.AddEntryEvent
	if !result.Success {
		msg := "unknown error"
		if result.Message != nil && *result.Message != "" {
			msg = *result.Message
		}
		return fmt.Errorf("server rejected event directive: %s", msg)
	}
	fmt.Fprintf(cmd.OutOrStdout(), "Event directive added for type %q in %s on %s\n",
		event.Type, eventLedger, event.Date)
	return nil
}
