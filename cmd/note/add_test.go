package note

import (
	"strings"
	"testing"
	"time"

	"beancount.io/beancount-cli/generated"
)

func TestBuildNoteInput(t *testing.T) {
	today := time.Now().Format("2006-01-02")

	tests := []struct {
		name    string
		date    string
		account string
		content string
		check   func(t *testing.T, got generated.LedgerNoteInput)
		wantErr string
	}{
		{
			name:    "all fields set",
			date:    "2024-01-15",
			account: "Assets:Checking",
			content: "Transferred funds from savings",
			check: func(t *testing.T, got generated.LedgerNoteInput) {
				t.Helper()
				if got.Date != "2024-01-15" {
					t.Errorf("Date: got %q, want %q", got.Date, "2024-01-15")
				}
				if got.Account != "Assets:Checking" {
					t.Errorf("Account: got %q, want %q", got.Account, "Assets:Checking")
				}
				if got.Content != "Transferred funds from savings" {
					t.Errorf("Content: got %q, want %q", got.Content, "Transferred funds from savings")
				}
			},
		},
		{
			name:    "empty date defaults to today",
			date:    "",
			account: "Assets:Cash",
			content: "Some note",
			check: func(t *testing.T, got generated.LedgerNoteInput) {
				t.Helper()
				if got.Date != today {
					t.Errorf("Date: got %q, want %q (today)", got.Date, today)
				}
			},
		},
		{
			name:    "content with special characters",
			date:    "2024-06-01",
			account: "Expenses:Food",
			content: "Lunch at Joe's Diner — $12.50",
			check: func(t *testing.T, got generated.LedgerNoteInput) {
				t.Helper()
				if got.Content != "Lunch at Joe's Diner — $12.50" {
					t.Errorf("Content: got %q, want %q", got.Content, "Lunch at Joe's Diner — $12.50")
				}
			},
		},
		{
			name:    "account is required",
			date:    "2024-01-15",
			account: "",
			content: "Some note",
			wantErr: "account is required",
		},
		{
			name:    "content is required",
			date:    "2024-01-15",
			account: "Assets:Checking",
			content: "",
			wantErr: "content is required",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := buildNoteInput(tc.date, tc.account, tc.content)
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
