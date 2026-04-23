package balance

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"beancount.io/beancount-cli/generated"
	"beancount.io/beancount-cli/internal/utils"
)

func NewCmdAdd() *cobra.Command {
	var ledger, account, date, amount, currency string

	cmd := &cobra.Command{
		Use:   "add",
		Short: "Add a balance directive to a ledger",
		Long: `Add a balance assertion directive to a Beancount ledger.

Example:
  beancount-cli balance add \
    --ledger user/mybook \
    --account Assets:Checking \
    --amount 1000.00 \
    --currency USD \
    --date 2024-01-01`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			balance, err := buildBalanceInput(date, account, amount, currency)
			if err != nil {
				return err
			}
			client, err := utils.NewAuthedClient()
			if err != nil {
				return err
			}
			resp, err := generated.AddEntryBalance(context.Background(), client, utils.LedgerID(ledger), balance)
			if err != nil {
				return fmt.Errorf("failed to add balance directive: %w", err)
			}
			result := resp.AddEntryBalance
			if !result.Success {
				msg := "unknown error"
				if result.Message != nil && *result.Message != "" {
					msg = *result.Message
				}
				return fmt.Errorf("server rejected balance directive: %s", msg)
			}
			fmt.Fprintf(cmd.OutOrStdout(), "Balance directive added for %s in %s on %s: %s %s\n",
				balance.Account, ledger, balance.Date, balance.Amount.Number, balance.Amount.Currency)
			return nil
		},
	}

	cmd.Flags().StringVar(&ledger, "ledger", "", "Ledger full name (e.g. user/mybook) (required)")
	cmd.Flags().StringVar(&account, "account", "", "Account name to assert (e.g. Assets:Checking) (required)")
	cmd.Flags().StringVar(&amount, "amount", "", "Balance amount (e.g. 1000.00) (required)")
	cmd.Flags().StringVar(&currency, "currency", "", "Currency code (e.g. USD) (required)")
	cmd.Flags().StringVar(&date, "date", "", "Balance date in YYYY-MM-DD format (default: today)")
	_ = cmd.MarkFlagRequired("ledger")
	_ = cmd.MarkFlagRequired("account")
	_ = cmd.MarkFlagRequired("amount")
	_ = cmd.MarkFlagRequired("currency")
	return cmd
}

func buildBalanceInput(date, account, amount, currency string) (generated.LedgerBalanceInput, error) {
	if account == "" {
		return generated.LedgerBalanceInput{}, fmt.Errorf("account is required")
	}
	if amount == "" {
		return generated.LedgerBalanceInput{}, fmt.Errorf("amount is required")
	}
	if currency == "" {
		return generated.LedgerBalanceInput{}, fmt.Errorf("currency is required")
	}
	currency = strings.ToUpper(strings.TrimSpace(currency))
	if date == "" {
		date = time.Now().Format("2006-01-02")
	}
	return generated.LedgerBalanceInput{
		Date:    date,
		Account: account,
		Amount: generated.LedgerAmountInput{
			Number:   amount,
			Currency: currency,
		},
	}, nil
}
