package collaborator

import (
	"strings"
	"testing"

	"beancount.io/beancount-cli/internal/utils"
)

func TestBuildAddInput(t *testing.T) {
	tests := []struct {
		name       string
		ledger     string
		username   string
		permission string
		checkID    func(t *testing.T, ledgerID string)
		wantCollab string
		wantPerm   string
	}{
		{
			name:       "default permission",
			ledger:     "user/mybook",
			username:   "alice",
			permission: "",
			checkID: func(t *testing.T, ledgerID string) {
				if ledgerID != utils.LedgerID("user/mybook") {
					t.Errorf("ledgerID = %q, want LedgerID(user/mybook)", ledgerID)
				}
			},
			wantCollab: "alice",
			wantPerm:   "read",
		},
		{
			name:       "explicit write permission",
			ledger:     "user/mybook",
			username:   "bob",
			permission: "write",
			checkID: func(t *testing.T, ledgerID string) {
				if !strings.Contains(ledgerID, "=") || ledgerID == "" {
					t.Errorf("ledgerID looks unencoded: %q", ledgerID)
				}
			},
			wantCollab: "bob",
			wantPerm:   "write",
		},
		{
			name:       "explicit read permission preserved",
			ledger:     "org/ledger",
			username:   "carol",
			permission: "read",
			checkID:    func(t *testing.T, ledgerID string) {},
			wantCollab: "carol",
			wantPerm:   "read",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ledgerID, collaborator, perm := buildAddInput(tc.ledger, tc.username, tc.permission)
			tc.checkID(t, ledgerID)
			if collaborator != tc.wantCollab {
				t.Errorf("collaborator = %q, want %q", collaborator, tc.wantCollab)
			}
			if perm != tc.wantPerm {
				t.Errorf("permission = %q, want %q", perm, tc.wantPerm)
			}
		})
	}
}
