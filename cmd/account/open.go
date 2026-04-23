package account

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"beancount.io/beancount-cli/generated"
	"beancount.io/beancount-cli/internal/utils"
)

func NewCmdOpen() *cobra.Command {
	var ledger, account, date, currencies string

	cmd := &cobra.Command{
		Use:   "open",
		Short: "Open an account in a ledger",
		Long: `Add an account open directive to a Beancount ledger.

Example:
  beancount-cli account open \
    --ledger user/mybook \
    --account Expenses:Food:Coffee \
    --date 2024-01-01 \
    --currency USD,EUR`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			open, err := buildOpenInput(date, account, currencies)
			if err != nil {
				return err
			}
			client, err := utils.NewAuthedClient()
			if err != nil {
				return err
			}
			resp, err := generated.AddEntryOpen(context.Background(), client, utils.LedgerID(ledger), open)
			if err != nil {
				return fmt.Errorf("failed to open account: %w", err)
			}
			result := resp.AddEntryOpen
			if !result.Success {
				msg := "unknown error"
				if result.Message != nil && *result.Message != "" {
					msg = *result.Message
				}
				return fmt.Errorf("server rejected open directive: %s", msg)
			}
			fmt.Fprintf(cmd.OutOrStdout(), "Account %s opened in %s on %s\n", open.Account, ledger, open.Date)
			return nil
		},
	}

	cmd.Flags().StringVar(&ledger, "ledger", "", "Ledger full name (e.g. user/mybook) (required)")
	cmd.Flags().StringVar(&account, "account", "", "Account name to open (e.g. Expenses:Food) (required)")
	cmd.Flags().StringVar(&date, "date", "", "Open date in YYYY-MM-DD format (default: today)")
	cmd.Flags().StringVar(&currencies, "currency", "", "Allowed currencies, comma-separated (e.g. USD,EUR)")
	_ = cmd.MarkFlagRequired("ledger")
	_ = cmd.MarkFlagRequired("account")
	return cmd
}

func buildOpenInput(date, account, currencies string) (generated.LedgerOpenInput, error) {
	if account == "" {
		return generated.LedgerOpenInput{}, fmt.Errorf("account is required")
	}
	if date == "" {
		date = time.Now().Format("2006-01-02")
	}
	parsed := []string{}
	for _, c := range strings.Split(currencies, ",") {
		if c = strings.TrimSpace(c); c != "" {
			parsed = append(parsed, c)
		}
	}
	return generated.LedgerOpenInput{
		Date:       date,
		Account:    account,
		Currencies: parsed,
	}, nil
}
