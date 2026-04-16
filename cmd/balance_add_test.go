package cmd

import (
	"strings"
	"testing"
	"time"

	"beancount.io/beancount-cli/generated"
)

func TestBuildBalanceInput(t *testing.T) {
	today := time.Now().Format("2006-01-02")

	tests := []struct {
		name     string
		date     string
		account  string
		amount   string
		currency string
		check    func(t *testing.T, got generated.LedgerBalanceInput)
		wantErr  string
	}{
		{
			name:     "all fields set",
			date:     "2024-01-15",
			account:  "Assets:Checking",
			amount:   "1000.00",
			currency: "USD",
			check: func(t *testing.T, got generated.LedgerBalanceInput) {
				t.Helper()
				if got.Date != "2024-01-15" {
					t.Errorf("Date: got %q, want %q", got.Date, "2024-01-15")
				}
				if got.Account != "Assets:Checking" {
					t.Errorf("Account: got %q, want %q", got.Account, "Assets:Checking")
				}
				if got.Amount.Number != "1000.00" {
					t.Errorf("Amount.Number: got %q, want %q", got.Amount.Number, "1000.00")
				}
				if got.Amount.Currency != "USD" {
					t.Errorf("Amount.Currency: got %q, want %q", got.Amount.Currency, "USD")
				}
			},
		},
		{
			name:     "empty date defaults to today",
			date:     "",
			account:  "Assets:Cash",
			amount:   "500.00",
			currency: "EUR",
			check: func(t *testing.T, got generated.LedgerBalanceInput) {
				t.Helper()
				if got.Date != today {
					t.Errorf("Date: got %q, want %q (today)", got.Date, today)
				}
			},
		},
		{
			name:     "currency is uppercased",
			date:     "2024-01-15",
			account:  "Assets:Cash",
			amount:   "100.00",
			currency: "usd",
			check: func(t *testing.T, got generated.LedgerBalanceInput) {
				t.Helper()
				if got.Amount.Currency != "USD" {
					t.Errorf("Amount.Currency: got %q, want %q", got.Amount.Currency, "USD")
				}
			},
		},
		{
			name:     "negative amount",
			date:     "2024-01-15",
			account:  "Liabilities:CreditCard",
			amount:   "-250.50",
			currency: "USD",
			check: func(t *testing.T, got generated.LedgerBalanceInput) {
				t.Helper()
				if got.Amount.Number != "-250.50" {
					t.Errorf("Amount.Number: got %q, want %q", got.Amount.Number, "-250.50")
				}
			},
		},
		{
			name:     "account is required",
			date:     "2024-01-15",
			account:  "",
			amount:   "100.00",
			currency: "USD",
			wantErr:  "account is required",
		},
		{
			name:     "amount is required",
			date:     "2024-01-15",
			account:  "Assets:Checking",
			amount:   "",
			currency: "USD",
			wantErr:  "amount is required",
		},
		{
			name:     "currency is required",
			date:     "2024-01-15",
			account:  "Assets:Checking",
			amount:   "100.00",
			currency: "",
			wantErr:  "currency is required",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := buildBalanceInput(tc.date, tc.account, tc.amount, tc.currency)
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
