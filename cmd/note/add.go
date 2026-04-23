package note

import (
	"context"
	"fmt"
	"time"

	"github.com/spf13/cobra"

	"beancount.io/beancount-cli/generated"
	"beancount.io/beancount-cli/internal/utils"
)

func NewCmdAdd() *cobra.Command {
	var ledger, account, date, content string

	cmd := &cobra.Command{
		Use:   "add",
		Short: "Add a note directive to a ledger",
		Long: `Add a note directive to a Beancount ledger.

Example:
  beancount-cli note add \
    --ledger user/mybook \
    --account Assets:Checking \
    --content "Transferred funds from savings" \
    --date 2024-01-01`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			note, err := buildNoteInput(date, account, content)
			if err != nil {
				return err
			}
			client, err := utils.NewAuthedClient()
			if err != nil {
				return err
			}
			resp, err := generated.AddEntryNote(context.Background(), client, utils.LedgerID(ledger), note)
			if err != nil {
				return fmt.Errorf("failed to add note directive: %w", err)
			}
			result := resp.AddEntryNote
			if !result.Success {
				msg := "unknown error"
				if result.Message != nil && *result.Message != "" {
					msg = *result.Message
				}
				return fmt.Errorf("server rejected note directive: %s", msg)
			}
			fmt.Fprintf(cmd.OutOrStdout(), "Note directive added for %s in %s on %s\n",
				note.Account, ledger, note.Date)
			return nil
		},
	}

	cmd.Flags().StringVar(&ledger, "ledger", "", "Ledger full name (e.g. user/mybook) (required)")
	cmd.Flags().StringVar(&account, "account", "", "Account name to attach the note to (required)")
	cmd.Flags().StringVar(&content, "content", "", "Note content (required)")
	cmd.Flags().StringVar(&date, "date", "", "Note date in YYYY-MM-DD format (default: today)")
	_ = cmd.MarkFlagRequired("ledger")
	_ = cmd.MarkFlagRequired("account")
	_ = cmd.MarkFlagRequired("content")
	return cmd
}

func buildNoteInput(date, account, content string) (generated.LedgerNoteInput, error) {
	if account == "" {
		return generated.LedgerNoteInput{}, fmt.Errorf("account is required")
	}
	if content == "" {
		return generated.LedgerNoteInput{}, fmt.Errorf("content is required")
	}
	if date == "" {
		date = time.Now().Format("2006-01-02")
	}
	return generated.LedgerNoteInput{
		Date:    date,
		Account: account,
		Content: content,
	}, nil
}
