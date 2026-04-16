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

var balanceAddCmd = &cobra.Command{
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
	RunE: runBalanceAdd,
}

var (
	balanceLedger   string
	balanceAccount  string
	balanceDate     string
	balanceAmount   string
	balanceCurrency string
)

func init() {
	balanceCmd.AddCommand(balanceAddCmd)
	balanceAddCmd.Flags().StringVar(&balanceLedger, "ledger", "", "Ledger full name (e.g. user/mybook) (required)")
	balanceAddCmd.Flags().StringVar(&balanceAccount, "account", "", "Account name to assert (e.g. Assets:Checking) (required)")
	balanceAddCmd.Flags().StringVar(&balanceAmount, "amount", "", "Balance amount (e.g. 1000.00) (required)")
	balanceAddCmd.Flags().StringVar(&balanceCurrency, "currency", "", "Currency code (e.g. USD) (required)")
	balanceAddCmd.Flags().StringVar(&balanceDate, "date", "", "Balance date in YYYY-MM-DD format (default: today)")
	_ = balanceAddCmd.MarkFlagRequired("ledger")
	_ = balanceAddCmd.MarkFlagRequired("account")
	_ = balanceAddCmd.MarkFlagRequired("amount")
	_ = balanceAddCmd.MarkFlagRequired("currency")
}

// buildBalanceInput constructs a LedgerBalanceInput from raw CLI values.
// Pure function — no I/O, no network calls.
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

func runBalanceAdd(cmd *cobra.Command, _ []string) error {
	balance, err := buildBalanceInput(balanceDate, balanceAccount, balanceAmount, balanceCurrency)
	if err != nil {
		return err
	}

	client, err := utils.NewAuthedClient()
	if err != nil {
		return err
	}

	resp, err := generated.AddEntryBalance(context.Background(), client, utils.LedgerID(balanceLedger), balance)
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
		balance.Account, balanceLedger, balance.Date, balance.Amount.Number, balance.Amount.Currency)
	return nil
}
