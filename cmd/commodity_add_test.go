package cmd

import (
	"strings"
	"testing"
	"time"

	"beancount.io/beancount-cli/generated"
)

func TestBuildCommodityInput(t *testing.T) {
	today := time.Now().Format("2006-01-02")

	tests := []struct {
		name     string
		date     string
		currency string
		check    func(t *testing.T, got generated.LedgerCommodityInput)
		wantErr  string
	}{
		{
			name:     "all fields set",
			date:     "2024-01-15",
			currency: "AAPL",
			check: func(t *testing.T, got generated.LedgerCommodityInput) {
				t.Helper()
				if got.Date != "2024-01-15" {
					t.Errorf("Date: got %q, want %q", got.Date, "2024-01-15")
				}
				if got.Currency != "AAPL" {
					t.Errorf("Currency: got %q, want %q", got.Currency, "AAPL")
				}
			},
		},
		{
			name:     "empty date defaults to today",
			date:     "",
			currency: "GOOG",
			check: func(t *testing.T, got generated.LedgerCommodityInput) {
				t.Helper()
				if got.Date != today {
					t.Errorf("Date: got %q, want %q (today)", got.Date, today)
				}
			},
		},
		{
			name:     "currency is uppercased",
			date:     "2024-01-15",
			currency: "aapl",
			check: func(t *testing.T, got generated.LedgerCommodityInput) {
				t.Helper()
				if got.Currency != "AAPL" {
					t.Errorf("Currency: got %q, want %q", got.Currency, "AAPL")
				}
			},
		},
		{
			name:     "currency is required",
			date:     "2024-01-15",
			currency: "",
			wantErr:  "currency is required",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := buildCommodityInput(tc.date, tc.currency)
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
