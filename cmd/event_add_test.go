package cmd

import (
	"strings"
	"testing"
	"time"

	"beancount.io/beancount-cli/generated"
)

func TestBuildEventInput(t *testing.T) {
	today := time.Now().Format("2006-01-02")

	tests := []struct {
		name        string
		date        string
		eventType   string
		description string
		check       func(t *testing.T, got generated.LedgerEventInput)
		wantErr     string
	}{
		{
			name:        "all fields set",
			date:        "2024-01-15",
			eventType:   "location",
			description: "New York",
			check: func(t *testing.T, got generated.LedgerEventInput) {
				t.Helper()
				if got.Date != "2024-01-15" {
					t.Errorf("Date: got %q, want %q", got.Date, "2024-01-15")
				}
				if got.Type != "location" {
					t.Errorf("Type: got %q, want %q", got.Type, "location")
				}
				if got.Description != "New York" {
					t.Errorf("Description: got %q, want %q", got.Description, "New York")
				}
			},
		},
		{
			name:        "empty date defaults to today",
			date:        "",
			eventType:   "location",
			description: "London",
			check: func(t *testing.T, got generated.LedgerEventInput) {
				t.Helper()
				if got.Date != today {
					t.Errorf("Date: got %q, want %q (today)", got.Date, today)
				}
			},
		},
		{
			name:        "description with special characters",
			date:        "2024-06-01",
			eventType:   "note",
			description: "Arrived at Joe's place — 9pm",
			check: func(t *testing.T, got generated.LedgerEventInput) {
				t.Helper()
				if got.Description != "Arrived at Joe's place — 9pm" {
					t.Errorf("Description: got %q, want %q", got.Description, "Arrived at Joe's place — 9pm")
				}
			},
		},
		{
			name:        "type is required",
			date:        "2024-01-15",
			eventType:   "",
			description: "Some description",
			wantErr:     "type is required",
		},
		{
			name:        "description is required",
			date:        "2024-01-15",
			eventType:   "location",
			description: "",
			wantErr:     "description is required",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := buildEventInput(tc.date, tc.eventType, tc.description)
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
