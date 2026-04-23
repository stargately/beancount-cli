package price

import (
	"strings"
	"testing"
	"time"

	"beancount.io/beancount-cli/generated"
)

func TestParseAmountFlag(t *testing.T) {
	tests := []struct {
		name         string
		input        string
		wantNumber   string
		wantCurrency string
		wantErr      string
	}{
		{
			name:         "valid integer amount",
			input:        "60000,USD",
			wantNumber:   "60000",
			wantCurrency: "USD",
		},
		{
			name:         "valid decimal amount",
			input:        "1234.56,EUR",
			wantNumber:   "1234.56",
			wantCurrency: "EUR",
		},
		{
			name:         "currency is uppercased",
			input:        "100,usd",
			wantNumber:   "100",
			wantCurrency: "USD",
		},
		{
			name:    "missing currency",
			input:   "100,",
			wantErr: "invalid amount",
		},
		{
			name:    "missing number",
			input:   ",USD",
			wantErr: "invalid amount",
		},
		{
			name:    "no comma",
			input:   "100USD",
			wantErr: "invalid amount",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			num, cur, err := parseAmountFlag(tc.input)
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
			if num != tc.wantNumber {
				t.Errorf("number: got %q, want %q", num, tc.wantNumber)
			}
			if cur != tc.wantCurrency {
				t.Errorf("currency: got %q, want %q", cur, tc.wantCurrency)
			}
		})
	}
}

func TestBuildPriceInput(t *testing.T) {
	today := time.Now().Format("2006-01-02")

	tests := []struct {
		name     string
		date     string
		currency string
		amount   string
		check    func(t *testing.T, got generated.LedgerPriceInput)
		wantErr  string
	}{
		{
			name:     "all fields set",
			date:     "2024-01-15",
			currency: "BTC",
			amount:   "60000,USD",
			check: func(t *testing.T, got generated.LedgerPriceInput) {
				t.Helper()
				if got.Date != "2024-01-15" {
					t.Errorf("Date: got %q, want %q", got.Date, "2024-01-15")
				}
				if got.Currency != "BTC" {
					t.Errorf("Currency: got %q, want %q", got.Currency, "BTC")
				}
				if got.Amount.Number != "60000" {
					t.Errorf("Amount.Number: got %q, want %q", got.Amount.Number, "60000")
				}
				if got.Amount.Currency != "USD" {
					t.Errorf("Amount.Currency: got %q, want %q", got.Amount.Currency, "USD")
				}
			},
		},
		{
			name:     "empty date defaults to today",
			date:     "",
			currency: "ETH",
			amount:   "3000,USD",
			check: func(t *testing.T, got generated.LedgerPriceInput) {
				t.Helper()
				if got.Date != today {
					t.Errorf("Date: got %q, want %q (today)", got.Date, today)
				}
			},
		},
		{
			name:     "currency is uppercased",
			date:     "2024-01-15",
			currency: "btc",
			amount:   "60000,usd",
			check: func(t *testing.T, got generated.LedgerPriceInput) {
				t.Helper()
				if got.Currency != "BTC" {
					t.Errorf("Currency: got %q, want %q", got.Currency, "BTC")
				}
				if got.Amount.Currency != "USD" {
					t.Errorf("Amount.Currency: got %q, want %q", got.Amount.Currency, "USD")
				}
			},
		},
		{
			name:     "currency is required",
			date:     "2024-01-15",
			currency: "",
			amount:   "60000,USD",
			wantErr:  "currency is required",
		},
		{
			name:     "amount is required",
			date:     "2024-01-15",
			currency: "BTC",
			amount:   "",
			wantErr:  "amount is required",
		},
		{
			name:     "invalid amount format",
			date:     "2024-01-15",
			currency: "BTC",
			amount:   "60000USD",
			wantErr:  "invalid amount",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := buildPriceInput(tc.date, tc.currency, tc.amount)
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
