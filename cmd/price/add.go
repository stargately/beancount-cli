package price

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
	var ledger, currency, amount, date string

	cmd := &cobra.Command{
		Use:   "add",
		Short: "Add a price directive to a ledger",
		Long: `Add a price directive to a Beancount ledger.

Example:
  beancount-cli price add \
    --ledger user/mybook \
    --currency BTC \
    --amount 60000,USD \
    --date 2024-01-01`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			priceInput, err := buildPriceInput(date, currency, amount)
			if err != nil {
				return err
			}
			client, err := utils.NewAuthedClient()
			if err != nil {
				return err
			}
			resp, err := generated.AddEntryPrice(context.Background(), client, utils.LedgerID(ledger), priceInput)
			if err != nil {
				return fmt.Errorf("failed to add price directive: %w", err)
			}
			result := resp.AddEntryPrice
			if !result.Success {
				msg := "unknown error"
				if result.Message != nil && *result.Message != "" {
					msg = *result.Message
				}
				return fmt.Errorf("server rejected price directive: %s", msg)
			}
			fmt.Fprintf(cmd.OutOrStdout(), "Price directive added for %s in %s on %s: %s %s\n",
				priceInput.Currency, ledger, priceInput.Date, priceInput.Amount.Number, priceInput.Amount.Currency)
			return nil
		},
	}

	cmd.Flags().StringVar(&ledger, "ledger", "", "Ledger full name (e.g. user/mybook) (required)")
	cmd.Flags().StringVar(&currency, "currency", "", "Commodity currency being priced (e.g. BTC) (required)")
	cmd.Flags().StringVar(&amount, "amount", "", "Price amount as number,currency (e.g. 60000,USD) (required)")
	cmd.Flags().StringVar(&date, "date", "", "Price date in YYYY-MM-DD format (default: today)")
	_ = cmd.MarkFlagRequired("ledger")
	_ = cmd.MarkFlagRequired("currency")
	_ = cmd.MarkFlagRequired("amount")
	return cmd
}

func parseAmountFlag(s string) (string, string, error) {
	return utils.ParseAmountFlag(s)
}

func buildPriceInput(date, currency, amount string) (generated.LedgerPriceInput, error) {
	if currency == "" {
		return generated.LedgerPriceInput{}, fmt.Errorf("currency is required")
	}
	if amount == "" {
		return generated.LedgerPriceInput{}, fmt.Errorf("amount is required")
	}
	number, amountCurrency, err := utils.ParseAmountFlag(amount)
	if err != nil {
		return generated.LedgerPriceInput{}, err
	}
	if date == "" {
		date = time.Now().Format("2006-01-02")
	}
	return generated.LedgerPriceInput{
		Date:     date,
		Currency: strings.ToUpper(strings.TrimSpace(currency)),
		Amount: generated.LedgerAmountInput{
			Number:   number,
			Currency: amountCurrency,
		},
	}, nil
}
