package document

import (
	"strings"
	"testing"
	"time"

	"beancount.io/beancount-cli/generated"
)

func TestBuildDocumentInput(t *testing.T) {
	today := time.Now().Format("2006-01-02")

	tests := []struct {
		name     string
		date     string
		account  string
		filename string
		tags     []string
		links    []string
		check    func(t *testing.T, got generated.LedgerDocumentInput)
		wantErr  string
	}{
		{
			name:     "all fields set",
			date:     "2024-01-15",
			account:  "Assets:Bank",
			filename: "receipts/2024-01-01.pdf",
			tags:     []string{"tax", "receipt"},
			links:    []string{"link1"},
			check: func(t *testing.T, got generated.LedgerDocumentInput) {
				t.Helper()
				if got.Date != "2024-01-15" {
					t.Errorf("Date: got %q, want %q", got.Date, "2024-01-15")
				}
				if got.Account != "Assets:Bank" {
					t.Errorf("Account: got %q, want %q", got.Account, "Assets:Bank")
				}
				if got.Filename != "receipts/2024-01-01.pdf" {
					t.Errorf("Filename: got %q, want %q", got.Filename, "receipts/2024-01-01.pdf")
				}
				if len(got.Tags) != 2 || got.Tags[0] != "tax" || got.Tags[1] != "receipt" {
					t.Errorf("Tags: got %v, want [tax receipt]", got.Tags)
				}
				if len(got.Links) != 1 || got.Links[0] != "link1" {
					t.Errorf("Links: got %v, want [link1]", got.Links)
				}
			},
		},
		{
			name:     "empty date defaults to today",
			date:     "",
			account:  "Assets:Bank",
			filename: "file.pdf",
			check: func(t *testing.T, got generated.LedgerDocumentInput) {
				t.Helper()
				if got.Date != today {
					t.Errorf("Date: got %q, want %q (today)", got.Date, today)
				}
			},
		},
		{
			name:     "nil tags defaults to empty slice",
			date:     "2024-01-15",
			account:  "Assets:Bank",
			filename: "file.pdf",
			tags:     nil,
			check: func(t *testing.T, got generated.LedgerDocumentInput) {
				t.Helper()
				if got.Tags == nil {
					t.Errorf("Tags: got nil, want []string{}")
				}
				if len(got.Tags) != 0 {
					t.Errorf("Tags: got %v, want empty slice", got.Tags)
				}
			},
		},
		{
			name:     "nil links defaults to empty slice",
			date:     "2024-01-15",
			account:  "Assets:Bank",
			filename: "file.pdf",
			links:    nil,
			check: func(t *testing.T, got generated.LedgerDocumentInput) {
				t.Helper()
				if got.Links == nil {
					t.Errorf("Links: got nil, want []string{}")
				}
				if len(got.Links) != 0 {
					t.Errorf("Links: got %v, want empty slice", got.Links)
				}
			},
		},
		{
			name:     "account is required",
			date:     "2024-01-15",
			account:  "",
			filename: "file.pdf",
			wantErr:  "account is required",
		},
		{
			name:     "filename is required",
			date:     "2024-01-15",
			account:  "Assets:Bank",
			filename: "",
			wantErr:  "filename is required",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := buildDocumentInput(tc.date, tc.account, tc.filename, tc.tags, tc.links)
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
