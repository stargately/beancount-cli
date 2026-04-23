package transaction

import (
	"encoding/base64"
	"strings"
	"testing"
	"time"

	"beancount.io/beancount-cli/generated"
	"beancount.io/beancount-cli/internal/utils"
)

// ----------------------------------------------------------------------------
// TestParsePosting
// ----------------------------------------------------------------------------

func TestParsePosting(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    generated.LedgerPostingInput
		wantErr string
	}{
		{
			name:  "valid simple",
			input: "Assets:Checking,100.00,USD",
			want: generated.LedgerPostingInput{
				Account: "Assets:Checking",
				Units:   generated.LedgerAmountInput{Number: "100.00", Currency: "USD"},
			},
		},
		{
			name:  "valid nested account",
			input: "Expenses:Food:Coffee,-5.50,EUR",
			want: generated.LedgerPostingInput{
				Account: "Expenses:Food:Coffee",
				Units:   generated.LedgerAmountInput{Number: "-5.50", Currency: "EUR"},
			},
		},
		{
			name:  "valid crypto",
			input: "Assets:Wallet,0.001,BTC",
			want: generated.LedgerPostingInput{
				Account: "Assets:Wallet",
				Units:   generated.LedgerAmountInput{Number: "0.001", Currency: "BTC"},
			},
		},
		{
			name:    "missing currency (2 parts)",
			input:   "Assets:Checking,100.00",
			wantErr: "invalid posting",
		},
		{
			name:    "empty string",
			input:   "",
			wantErr: "invalid posting",
		},
		{
			name:    "single part only",
			input:   "Assets:Checking",
			wantErr: "invalid posting",
		},
		{
			// SplitN(s, ",", 3) stops at 3 parts — extra commas fold into currency.
			// Documents the edge-case behavior rather than testing an error path.
			name:  "four comma-separated values (SplitN 3 behavior)",
			input: "a,b,c,d",
			want: generated.LedgerPostingInput{
				Account: "a",
				Units:   generated.LedgerAmountInput{Number: "b", Currency: "c,d"},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := parsePosting(tc.input)
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
			if got != tc.want {
				t.Errorf("got %+v, want %+v", got, tc.want)
			}
		})
	}
}

// ----------------------------------------------------------------------------
// TestBuildTransactionInput
// ----------------------------------------------------------------------------

func TestBuildTransactionInput(t *testing.T) {
	today := time.Now().Format("2006-01-02")

	tests := []struct {
		name        string
		flag        string
		date        string
		payee       string
		narration   string
		postingStrs []string
		tags        []string
		links       []string
		check       func(t *testing.T, tx generated.LedgerTransactionInput)
		wantErr     string
	}{
		{
			name:        "all fields set",
			flag:        "!",
			date:        "2024-01-15",
			payee:       "Starbucks",
			narration:   "Coffee",
			postingStrs: []string{"Expenses:Food,5.00,USD"},
			tags:        []string{"travel"},
			links:       []string{"^trip"},
			check: func(t *testing.T, tx generated.LedgerTransactionInput) {
				t.Helper()
				if tx.Date != "2024-01-15" {
					t.Errorf("Date: got %q, want %q", tx.Date, "2024-01-15")
				}
				if tx.Flag != "!" {
					t.Errorf("Flag: got %q, want %q", tx.Flag, "!")
				}
				if tx.Payee == nil || *tx.Payee != "Starbucks" {
					t.Errorf("Payee: got %v, want ptr(Starbucks)", tx.Payee)
				}
				if tx.Narration == nil || *tx.Narration != "Coffee" {
					t.Errorf("Narration: got %v, want ptr(Coffee)", tx.Narration)
				}
				if len(tx.Tags) != 1 || tx.Tags[0] != "travel" {
					t.Errorf("Tags: got %v, want [travel]", tx.Tags)
				}
				if len(tx.Links) != 1 || tx.Links[0] != "^trip" {
					t.Errorf("Links: got %v, want [^trip]", tx.Links)
				}
			},
		},
		{
			name:        "payee empty becomes nil pointer",
			flag:        "*",
			date:        "2024-01-15",
			payee:       "",
			postingStrs: []string{"Assets:Cash,10,USD"},
			check: func(t *testing.T, tx generated.LedgerTransactionInput) {
				t.Helper()
				if tx.Payee != nil {
					t.Errorf("Payee: got %v, want nil", tx.Payee)
				}
			},
		},
		{
			name:        "narration empty becomes nil pointer",
			flag:        "*",
			date:        "2024-01-15",
			narration:   "",
			postingStrs: []string{"Assets:Cash,10,USD"},
			check: func(t *testing.T, tx generated.LedgerTransactionInput) {
				t.Helper()
				if tx.Narration != nil {
					t.Errorf("Narration: got %v, want nil", tx.Narration)
				}
			},
		},
		{
			name:        "nil tags become empty slice",
			flag:        "*",
			date:        "2024-01-15",
			postingStrs: []string{"Assets:Cash,10,USD"},
			tags:        nil,
			check: func(t *testing.T, tx generated.LedgerTransactionInput) {
				t.Helper()
				if tx.Tags == nil {
					t.Error("Tags: got nil, want empty slice (avoids GraphQL serialization issue)")
				}
				if len(tx.Tags) != 0 {
					t.Errorf("Tags: got %v, want []", tx.Tags)
				}
			},
		},
		{
			name:        "nil links become empty slice",
			flag:        "*",
			date:        "2024-01-15",
			postingStrs: []string{"Assets:Cash,10,USD"},
			links:       nil,
			check: func(t *testing.T, tx generated.LedgerTransactionInput) {
				t.Helper()
				if tx.Links == nil {
					t.Error("Links: got nil, want empty slice (avoids GraphQL serialization issue)")
				}
				if len(tx.Links) != 0 {
					t.Errorf("Links: got %v, want []", tx.Links)
				}
			},
		},
		{
			name:        "empty date defaults to today",
			flag:        "*",
			date:        "",
			postingStrs: []string{"Assets:Cash,10,USD"},
			check: func(t *testing.T, tx generated.LedgerTransactionInput) {
				t.Helper()
				if tx.Date != today {
					t.Errorf("Date: got %q, want %q (today)", tx.Date, today)
				}
			},
		},
		{
			name:        "default cleared flag",
			flag:        "*",
			date:        "2024-01-15",
			postingStrs: []string{"Assets:Cash,10,USD"},
			check: func(t *testing.T, tx generated.LedgerTransactionInput) {
				t.Helper()
				if tx.Flag != "*" {
					t.Errorf("Flag: got %q, want \"*\"", tx.Flag)
				}
			},
		},
		{
			name: "multiple postings preserve order",
			flag: "*",
			date: "2024-01-15",
			postingStrs: []string{
				"Expenses:Food,5.00,USD",
				"Assets:Cash,-5.00,USD",
			},
			check: func(t *testing.T, tx generated.LedgerTransactionInput) {
				t.Helper()
				if len(tx.Postings) != 2 {
					t.Fatalf("Postings: got %d, want 2", len(tx.Postings))
				}
				if tx.Postings[0].Account != "Expenses:Food" {
					t.Errorf("Postings[0].Account: got %q, want Expenses:Food", tx.Postings[0].Account)
				}
				if tx.Postings[1].Account != "Assets:Cash" {
					t.Errorf("Postings[1].Account: got %q, want Assets:Cash", tx.Postings[1].Account)
				}
			},
		},
		{
			name:        "no postings returns error",
			flag:        "*",
			date:        "2024-01-15",
			postingStrs: nil,
			wantErr:     "at least one --posting is required",
		},
		{
			name:        "invalid posting in list returns error",
			flag:        "*",
			date:        "2024-01-15",
			postingStrs: []string{"bad-format"},
			wantErr:     "invalid posting",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tx, err := buildTransactionInput(tc.flag, tc.date, tc.payee, tc.narration, tc.postingStrs, tc.tags, tc.links)
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
				tc.check(t, tx)
			}
		})
	}
}

// ----------------------------------------------------------------------------
// TestLedgerID
// ----------------------------------------------------------------------------

func TestLedgerID(t *testing.T) {
	tests := []struct {
		fullName string
	}{
		{"user/mybook"},
		{"alice/personal"},
		{"org/team-ledger"},
	}
	for _, tc := range tests {
		t.Run(tc.fullName, func(t *testing.T) {
			got := utils.LedgerID(tc.fullName)
			want := base64.URLEncoding.EncodeToString([]byte(tc.fullName))
			if got != want {
				t.Errorf("LedgerID(%q) = %q, want %q", tc.fullName, got, want)
			}
		})
	}
}
