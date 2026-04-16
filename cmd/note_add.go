package cmd

import (
	"context"
	"fmt"
	"time"

	"github.com/spf13/cobra"

	"beancount.io/beancount-cli/generated"
	"beancount.io/beancount-cli/internal/utils"
)

var noteAddCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a note directive to a ledger",
	Long: `Add a note directive to a Beancount ledger.

Example:
  beancount-cli note add \
    --ledger user/mybook \
    --account Assets:Checking \
    --content "Transferred funds from savings" \
    --date 2024-01-01`,
	RunE: runNoteAdd,
}

var (
	noteLedger  string
	noteAccount string
	noteDate    string
	noteContent string
)

func init() {
	noteCmd.AddCommand(noteAddCmd)
	noteAddCmd.Flags().StringVar(&noteLedger, "ledger", "", "Ledger full name (e.g. user/mybook) (required)")
	noteAddCmd.Flags().StringVar(&noteAccount, "account", "", "Account name to attach the note to (required)")
	noteAddCmd.Flags().StringVar(&noteContent, "content", "", "Note content (required)")
	noteAddCmd.Flags().StringVar(&noteDate, "date", "", "Note date in YYYY-MM-DD format (default: today)")
	_ = noteAddCmd.MarkFlagRequired("ledger")
	_ = noteAddCmd.MarkFlagRequired("account")
	_ = noteAddCmd.MarkFlagRequired("content")
}

// buildNoteInput constructs a LedgerNoteInput from raw CLI values.
// Pure function — no I/O, no network calls.
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

func runNoteAdd(cmd *cobra.Command, _ []string) error {
	note, err := buildNoteInput(noteDate, noteAccount, noteContent)
	if err != nil {
		return err
	}

	client, err := utils.NewAuthedClient()
	if err != nil {
		return err
	}

	resp, err := generated.AddEntryNote(context.Background(), client, utils.LedgerID(noteLedger), note)
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
		note.Account, noteLedger, note.Date)
	return nil
}
