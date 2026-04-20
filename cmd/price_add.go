package cmd

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"beancount.io/beancount-cli/generated"
	"beancount.io/beancount-cli/internal/utils"
)

var priceAddCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a price directive to a ledger",
	Long: `Add a price directive to a Beancount ledger.

Example:
  beancount-cli price add \
    --ledger user/mybook \
    --currency BTC \
    --amount 60000,USD \
    --date 2024-01-01`,
	RunE: runPriceAdd,
}

var (
	priceLedger   string
	priceCurrency string
	priceAmount   string
	priceDate     string
)

func init() {
	priceCmd.AddCommand(priceAddCmd)
	priceAddCmd.Flags().StringVar(&priceLedger, "ledger", "", "Ledger full name (e.g. user/mybook) (required)")
	priceAddCmd.Flags().StringVar(&priceCurrency, "currency", "", "Commodity currency being priced (e.g. BTC) (required)")
	priceAddCmd.Flags().StringVar(&priceAmount, "amount", "", "Price amount as number,currency (e.g. 60000,USD) (required)")
	priceAddCmd.Flags().StringVar(&priceDate, "date", "", "Price date in YYYY-MM-DD format (default: today)")
	_ = priceAddCmd.MarkFlagRequired("ledger")
	_ = priceAddCmd.MarkFlagRequired("currency")
	_ = priceAddCmd.MarkFlagRequired("amount")
}

// parseAmountFlag parses a combined "number,currency" string (e.g. "60000,USD").
func parseAmountFlag(s string) (number, currency string, err error) {
	parts := strings.SplitN(s, ",", 2)
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return "", "", fmt.Errorf("invalid amount %q: expected number,currency (e.g. 60000,USD)", s)
	}
	return parts[0], strings.ToUpper(strings.TrimSpace(parts[1])), nil
}

func buildPriceInput(date, currency, amount string) (generated.LedgerPriceInput, error) {
	if currency == "" {
		return generated.LedgerPriceInput{}, fmt.Errorf("currency is required")
	}
	if amount == "" {
		return generated.LedgerPriceInput{}, fmt.Errorf("amount is required")
	}
	number, amountCurrency, err := parseAmountFlag(amount)
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

func runPriceAdd(cmd *cobra.Command, _ []string) error {
	price, err := buildPriceInput(priceDate, priceCurrency, priceAmount)
	if err != nil {
		return err
	}

	client, err := utils.NewAuthedClient()
	if err != nil {
		return err
	}

	resp, err := generated.AddEntryPrice(context.Background(), client, utils.LedgerID(priceLedger), price)
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
		price.Currency, priceLedger, price.Date, price.Amount.Number, price.Amount.Currency)
	return nil
}
