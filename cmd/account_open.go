package cmd

import (
	"context"
	"fmt"
	"time"

	"github.com/spf13/cobra"

	"beancount.io/beancount-cli/generated"
	"beancount.io/beancount-cli/internal/utils"
)

var accountOpenCmd = &cobra.Command{
	Use:   "open",
	Short: "Open an account in a ledger",
	Long: `Add an account open directive to a Beancount ledger.

Example:
  beancount-cli account open \
    --ledger user/mybook \
    --account Expenses:Food:Coffee \
    --date 2024-01-01 \
    --currency USD \
    --currency EUR`,
	RunE: runAccountOpen,
}

var (
	openLedger     string
	openAccount    string
	openDate       string
	openCurrencies []string
)

func init() {
	accountCmd.AddCommand(accountOpenCmd)
	accountOpenCmd.Flags().StringVar(&openLedger, "ledger", "", "Ledger full name (e.g. user/mybook) (required)")
	accountOpenCmd.Flags().StringVar(&openAccount, "account", "", "Account name to open (e.g. Expenses:Food) (required)")
	accountOpenCmd.Flags().StringVar(&openDate, "date", "", "Open date in YYYY-MM-DD format (default: today)")
	accountOpenCmd.Flags().StringArrayVar(&openCurrencies, "currency", nil, "Allowed currency (repeatable, optional)")
	_ = accountOpenCmd.MarkFlagRequired("ledger")
	_ = accountOpenCmd.MarkFlagRequired("account")
}

// buildOpenInput constructs a LedgerOpenInput from raw CLI values.
// Pure function — no I/O, no network calls.
func buildOpenInput(date, account string, currencies []string) (generated.LedgerOpenInput, error) {
	if account == "" {
		return generated.LedgerOpenInput{}, fmt.Errorf("account is required")
	}
	if date == "" {
		date = time.Now().Format("2006-01-02")
	}
	if currencies == nil {
		currencies = []string{}
	}
	return generated.LedgerOpenInput{
		Date:       date,
		Account:    account,
		Currencies: currencies,
	}, nil
}

func runAccountOpen(cmd *cobra.Command, _ []string) error {
	open, err := buildOpenInput(openDate, openAccount, openCurrencies)
	if err != nil {
		return err
	}

	client, err := utils.NewAuthedClient()
	if err != nil {
		return err
	}

	resp, err := generated.AddEntryOpen(context.Background(), client, utils.LedgerID(openLedger), open)
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

	fmt.Fprintf(cmd.OutOrStdout(), "Account %s opened in %s on %s\n", open.Account, openLedger, open.Date)
	return nil
}
