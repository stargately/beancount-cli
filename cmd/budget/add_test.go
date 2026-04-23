package budget

import (
	"strings"
	"testing"
	"time"

	"beancount.io/beancount-cli/generated"
)

func TestBuildBudgetInput(t *testing.T) {
	today := time.Now().Format("2006-01-02")

	tests := []struct {
		name     string
		date     string
		account  string
		amount   string
		interval string
		check    func(t *testing.T, got generated.LedgerBudgetInput)
		wantErr  string
	}{
		{
			name:     "all fields set",
			date:     "2024-01-15",
			account:  "Expenses:Food",
			amount:   "500,USD",
			interval: "MONTHLY",
			check: func(t *testing.T, got generated.LedgerBudgetInput) {
				t.Helper()
				if got.Date != "2024-01-15" {
					t.Errorf("Date: got %q, want %q", got.Date, "2024-01-15")
				}
				if got.Account != "Expenses:Food" {
					t.Errorf("Account: got %q, want %q", got.Account, "Expenses:Food")
				}
				if got.Amount.Number != "500" {
					t.Errorf("Amount.Number: got %q, want %q", got.Amount.Number, "500")
				}
				if got.Amount.Currency != "USD" {
					t.Errorf("Amount.Currency: got %q, want %q", got.Amount.Currency, "USD")
				}
				if got.Interval != generated.BudgetIntervalMonthly {
					t.Errorf("Interval: got %q, want %q", got.Interval, generated.BudgetIntervalMonthly)
				}
			},
		},
		{
			name:     "empty date defaults to today",
			date:     "",
			account:  "Expenses:Food",
			amount:   "500,USD",
			interval: "MONTHLY",
			check: func(t *testing.T, got generated.LedgerBudgetInput) {
				t.Helper()
				if got.Date != today {
					t.Errorf("Date: got %q, want %q (today)", got.Date, today)
				}
			},
		},
		{
			name:     "all valid intervals",
			date:     "2024-01-15",
			account:  "Expenses:Rent",
			amount:   "1200,USD",
			interval: "YEARLY",
			check: func(t *testing.T, got generated.LedgerBudgetInput) {
				t.Helper()
				if got.Interval != generated.BudgetIntervalYearly {
					t.Errorf("Interval: got %q, want %q", got.Interval, generated.BudgetIntervalYearly)
				}
			},
		},
		{
			name:     "interval case-insensitive",
			date:     "2024-01-15",
			account:  "Expenses:Food",
			amount:   "100,EUR",
			interval: "weekly",
			check: func(t *testing.T, got generated.LedgerBudgetInput) {
				t.Helper()
				if got.Interval != generated.BudgetIntervalWeekly {
					t.Errorf("Interval: got %q, want %q", got.Interval, generated.BudgetIntervalWeekly)
				}
			},
		},
		{
			name:     "account is required",
			date:     "2024-01-15",
			account:  "",
			amount:   "500,USD",
			interval: "MONTHLY",
			wantErr:  "account is required",
		},
		{
			name:     "amount is required",
			date:     "2024-01-15",
			account:  "Expenses:Food",
			amount:   "",
			interval: "MONTHLY",
			wantErr:  "amount is required",
		},
		{
			name:     "interval is required",
			date:     "2024-01-15",
			account:  "Expenses:Food",
			amount:   "500,USD",
			interval: "",
			wantErr:  "interval is required",
		},
		{
			name:     "invalid interval",
			date:     "2024-01-15",
			account:  "Expenses:Food",
			amount:   "500,USD",
			interval: "BIWEEKLY",
			wantErr:  "invalid interval",
		},
		{
			name:     "invalid amount format",
			date:     "2024-01-15",
			account:  "Expenses:Food",
			amount:   "500USD",
			interval: "MONTHLY",
			wantErr:  "invalid amount",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := buildBudgetInput(tc.date, tc.account, tc.amount, tc.interval)
			if tc.wantErr != "" {
				if err == nil {
					t.Fatalf("expected error containing %q, got nil", tc.wantErr)
				}
				if !strings.Contains(err.Error(), tc.wantErr) {
					t.Fatalf("error %q does not contain %q", err.Error(), tc.wantErr)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if tc.check != nil {
				tc.check(t, got)
			}
		})
	}
}
