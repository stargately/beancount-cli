package account

import (
	"strings"
	"testing"
	"time"

	"beancount.io/beancount-cli/generated"
)

func TestBuildCloseInput(t *testing.T) {
	today := time.Now().Format("2006-01-02")

	tests := []struct {
		name    string
		date    string
		account string
		check   func(t *testing.T, got generated.LedgerCloseInput)
		wantErr string
	}{
		{
			name:    "all fields set",
			date:    "2024-12-31",
			account: "Expenses:Food:Coffee",
			check: func(t *testing.T, got generated.LedgerCloseInput) {
				t.Helper()
				if got.Date != "2024-12-31" {
					t.Errorf("Date: got %q, want %q", got.Date, "2024-12-31")
				}
				if got.Account != "Expenses:Food:Coffee" {
					t.Errorf("Account: got %q, want %q", got.Account, "Expenses:Food:Coffee")
				}
			},
		},
		{
			name:    "empty date defaults to today",
			date:    "",
			account: "Assets:Checking",
			check: func(t *testing.T, got generated.LedgerCloseInput) {
				t.Helper()
				if got.Date != today {
					t.Errorf("Date: got %q, want %q (today)", got.Date, today)
				}
			},
		},
		{
			name:    "empty account returns error",
			date:    "2024-12-31",
			account: "",
			wantErr: "account",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := buildCloseInput(tc.date, tc.account)
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
