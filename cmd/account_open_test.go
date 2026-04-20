package cmd

import (
	"strings"
	"testing"
	"time"

	"beancount.io/beancount-cli/generated"
)

func TestBuildOpenInput(t *testing.T) {
	today := time.Now().Format("2006-01-02")

	tests := []struct {
		name       string
		date       string
		account    string
		currencies string
		check      func(t *testing.T, got generated.LedgerOpenInput)
		wantErr    string
	}{
		{
			name:       "all fields set",
			date:       "2024-01-01",
			account:    "Expenses:Food:Coffee",
			currencies: "USD,EUR",
			check: func(t *testing.T, got generated.LedgerOpenInput) {
				t.Helper()
				if got.Date != "2024-01-01" {
					t.Errorf("Date: got %q, want %q", got.Date, "2024-01-01")
				}
				if got.Account != "Expenses:Food:Coffee" {
					t.Errorf("Account: got %q, want %q", got.Account, "Expenses:Food:Coffee")
				}
				if len(got.Currencies) != 2 || got.Currencies[0] != "USD" || got.Currencies[1] != "EUR" {
					t.Errorf("Currencies: got %v, want [USD EUR]", got.Currencies)
				}
			},
		},
		{
			name:       "empty date defaults to today",
			date:       "",
			account:    "Assets:Checking",
			currencies: "",
			check: func(t *testing.T, got generated.LedgerOpenInput) {
				t.Helper()
				if got.Date != today {
					t.Errorf("Date: got %q, want %q (today)", got.Date, today)
				}
			},
		},
		{
			name:       "empty currencies become empty slice",
			date:       "2024-01-01",
			account:    "Assets:Checking",
			currencies: "",
			check: func(t *testing.T, got generated.LedgerOpenInput) {
				t.Helper()
				if got.Currencies == nil {
					t.Error("Currencies: got nil, want empty slice (avoids GraphQL serialization issue)")
				}
				if len(got.Currencies) != 0 {
					t.Errorf("Currencies: got %v, want []", got.Currencies)
				}
			},
		},
		{
			name:       "multiple currencies preserve order",
			date:       "2024-01-01",
			account:    "Assets:Wallet",
			currencies: "BTC,ETH,USD",
			check: func(t *testing.T, got generated.LedgerOpenInput) {
				t.Helper()
				if len(got.Currencies) != 3 {
					t.Fatalf("Currencies: got %d, want 3", len(got.Currencies))
				}
				if got.Currencies[0] != "BTC" {
					t.Errorf("Currencies[0]: got %q, want BTC", got.Currencies[0])
				}
				if got.Currencies[1] != "ETH" {
					t.Errorf("Currencies[1]: got %q, want ETH", got.Currencies[1])
				}
				if got.Currencies[2] != "USD" {
					t.Errorf("Currencies[2]: got %q, want USD", got.Currencies[2])
				}
			},
		},
		{
			name:       "currencies with spaces are trimmed",
			date:       "2024-01-01",
			account:    "Assets:Brokerage",
			currencies: "USD, EUR , GBP",
			check: func(t *testing.T, got generated.LedgerOpenInput) {
				t.Helper()
				if len(got.Currencies) != 3 {
					t.Fatalf("Currencies: got %d, want 3", len(got.Currencies))
				}
				if got.Currencies[1] != "EUR" {
					t.Errorf("Currencies[1]: got %q, want EUR", got.Currencies[1])
				}
			},
		},
		{
			name:    "empty account returns error",
			date:    "2024-01-01",
			account: "",
			wantErr: "account",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := buildOpenInput(tc.date, tc.account, tc.currencies)
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
