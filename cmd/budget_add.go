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

var budgetAddCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a budget directive to a ledger",
	Long: `Add a budget directive to a Beancount ledger.

Example:
  beancount-cli budget add \
    --ledger user/mybook \
    --account Expenses:Food \
    --amount 500,USD \
    --interval MONTHLY`,
	RunE: runBudgetAdd,
}

var (
	budgetLedger   string
	budgetAccount  string
	budgetAmount   string
	budgetInterval string
	budgetDate     string
)

var validBudgetIntervals = map[string]generated.BudgetInterval{
	"DAILY":     generated.BudgetIntervalDaily,
	"WEEKLY":    generated.BudgetIntervalWeekly,
	"MONTHLY":   generated.BudgetIntervalMonthly,
	"QUARTERLY": generated.BudgetIntervalQuarterly,
	"YEARLY":    generated.BudgetIntervalYearly,
}

func init() {
	budgetCmd.AddCommand(budgetAddCmd)
	budgetAddCmd.Flags().StringVar(&budgetLedger, "ledger", "", "Ledger full name (e.g. user/mybook) (required)")
	budgetAddCmd.Flags().StringVar(&budgetAccount, "account", "", "Account name (e.g. Expenses:Food) (required)")
	budgetAddCmd.Flags().StringVar(&budgetAmount, "amount", "", "Budget amount as number,currency (e.g. 500,USD) (required)")
	budgetAddCmd.Flags().StringVar(&budgetInterval, "interval", "", "Budget interval: DAILY, WEEKLY, MONTHLY, QUARTERLY, YEARLY (required)")
	budgetAddCmd.Flags().StringVar(&budgetDate, "date", "", "Budget date in YYYY-MM-DD format (default: today)")
	_ = budgetAddCmd.MarkFlagRequired("ledger")
	_ = budgetAddCmd.MarkFlagRequired("account")
	_ = budgetAddCmd.MarkFlagRequired("amount")
	_ = budgetAddCmd.MarkFlagRequired("interval")
}

func buildBudgetInput(date, account, amount, interval string) (generated.LedgerBudgetInput, error) {
	if account == "" {
		return generated.LedgerBudgetInput{}, fmt.Errorf("account is required")
	}
	if amount == "" {
		return generated.LedgerBudgetInput{}, fmt.Errorf("amount is required")
	}
	if interval == "" {
		return generated.LedgerBudgetInput{}, fmt.Errorf("interval is required")
	}
	budgetIntervalVal, ok := validBudgetIntervals[strings.ToUpper(strings.TrimSpace(interval))]
	if !ok {
		return generated.LedgerBudgetInput{}, fmt.Errorf("invalid interval %q: must be DAILY, WEEKLY, MONTHLY, QUARTERLY, or YEARLY", interval)
	}
	number, currency, err := parseAmountFlag(amount)
	if err != nil {
		return generated.LedgerBudgetInput{}, err
	}
	if date == "" {
		date = time.Now().Format("2006-01-02")
	}
	return generated.LedgerBudgetInput{
		Date:     date,
		Account:  account,
		Interval: budgetIntervalVal,
		Amount: generated.LedgerAmountInput{
			Number:   number,
			Currency: currency,
		},
	}, nil
}

func runBudgetAdd(cmd *cobra.Command, _ []string) error {
	budget, err := buildBudgetInput(budgetDate, budgetAccount, budgetAmount, budgetInterval)
	if err != nil {
		return err
	}

	client, err := utils.NewAuthedClient()
	if err != nil {
		return err
	}

	resp, err := generated.AddEntryBudget(context.Background(), client, utils.LedgerID(budgetLedger), budget)
	if err != nil {
		return fmt.Errorf("failed to add budget directive: %w", err)
	}

	result := resp.AddEntryBudget
	if !result.Success {
		msg := "unknown error"
		if result.Message != nil && *result.Message != "" {
			msg = *result.Message
		}
		return fmt.Errorf("server rejected budget directive: %s", msg)
	}

	fmt.Fprintf(cmd.OutOrStdout(), "Budget directive added for %s in %s on %s: %s %s %s\n",
		budget.Account, budgetLedger, budget.Date, budget.Amount.Number, budget.Amount.Currency, budget.Interval)
	return nil
}
