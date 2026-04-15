package cmd

import (
	"context"
	"fmt"
	"time"

	"github.com/spf13/cobra"

	"beancount.io/beancount-cli/generated"
	"beancount.io/beancount-cli/internal/utils"
)

var accountCloseCmd = &cobra.Command{
	Use:   "close",
	Short: "Close an account in a ledger",
	Long: `Add an account close directive to a Beancount ledger.

Example:
  beancount-cli account close \
    --ledger user/mybook \
    --account Expenses:Food:Coffee \
    --date 2024-12-31`,
	RunE: runAccountClose,
}

var (
	closeLedger  string
	closeAccount string
	closeDate    string
)

func init() {
	accountCmd.AddCommand(accountCloseCmd)
	accountCloseCmd.Flags().StringVar(&closeLedger, "ledger", "", "Ledger full name (e.g. user/mybook) (required)")
	accountCloseCmd.Flags().StringVar(&closeAccount, "account", "", "Account name to close (e.g. Expenses:Food) (required)")
	accountCloseCmd.Flags().StringVar(&closeDate, "date", "", "Close date in YYYY-MM-DD format (default: today)")
	_ = accountCloseCmd.MarkFlagRequired("ledger")
	_ = accountCloseCmd.MarkFlagRequired("account")
}

// buildCloseInput constructs a LedgerCloseInput from raw CLI values.
// Pure function — no I/O, no network calls.
func buildCloseInput(date, account string) (generated.LedgerCloseInput, error) {
	if account == "" {
		return generated.LedgerCloseInput{}, fmt.Errorf("account is required")
	}
	if date == "" {
		date = time.Now().Format("2006-01-02")
	}
	return generated.LedgerCloseInput{
		Date:    date,
		Account: account,
	}, nil
}

func runAccountClose(cmd *cobra.Command, _ []string) error {
	close, err := buildCloseInput(closeDate, closeAccount)
	if err != nil {
		return err
	}

	client, err := utils.NewAuthedClient()
	if err != nil {
		return err
	}

	resp, err := generated.AddEntryClose(context.Background(), client, utils.LedgerID(closeLedger), close)
	if err != nil {
		return fmt.Errorf("failed to close account: %w", err)
	}

	result := resp.AddEntryClose
	if !result.Success {
		msg := "unknown error"
		if result.Message != nil && *result.Message != "" {
			msg = *result.Message
		}
		return fmt.Errorf("server rejected close directive: %s", msg)
	}

	fmt.Fprintf(cmd.OutOrStdout(), "Account %s closed in %s on %s\n", close.Account, closeLedger, close.Date)
	return nil
}
