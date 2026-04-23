package account

import (
	"context"
	"fmt"
	"time"

	"github.com/spf13/cobra"

	"beancount.io/beancount-cli/generated"
	"beancount.io/beancount-cli/internal/utils"
)

func NewCmdClose() *cobra.Command {
	var ledger, account, date string

	cmd := &cobra.Command{
		Use:   "close",
		Short: "Close an account in a ledger",
		Long: `Add an account close directive to a Beancount ledger.

Example:
  beancount-cli account close \
    --ledger user/mybook \
    --account Expenses:Food:Coffee \
    --date 2024-12-31`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			close, err := buildCloseInput(date, account)
			if err != nil {
				return err
			}
			client, err := utils.NewAuthedClient()
			if err != nil {
				return err
			}
			resp, err := generated.AddEntryClose(context.Background(), client, utils.LedgerID(ledger), close)
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
			fmt.Fprintf(cmd.OutOrStdout(), "Account %s closed in %s on %s\n", close.Account, ledger, close.Date)
			return nil
		},
	}

	cmd.Flags().StringVar(&ledger, "ledger", "", "Ledger full name (e.g. user/mybook) (required)")
	cmd.Flags().StringVar(&account, "account", "", "Account name to close (e.g. Expenses:Food) (required)")
	cmd.Flags().StringVar(&date, "date", "", "Close date in YYYY-MM-DD format (default: today)")
	_ = cmd.MarkFlagRequired("ledger")
	_ = cmd.MarkFlagRequired("account")
	return cmd
}

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
