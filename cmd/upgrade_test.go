package cmd

import (
	"context"
	"strings"
	"testing"

	"beancount.io/beancount-cli/internal/config"
	"beancount.io/beancount-cli/internal/updater"
)

// ----------------------------------------------------------------------------
// TestCheckUpdate_DevBuild
// ----------------------------------------------------------------------------

func TestCheckUpdate_DevBuild(t *testing.T) {
	for _, ver := range []string{"dev", ""} {
		_, _, err := updater.CheckUpdate(context.Background(), ver, config.RepoSlug)
		if err == nil {
			t.Fatalf("version=%q: expected error, got nil", ver)
		}
		if !strings.Contains(err.Error(), "dev build") {
			t.Errorf("version=%q: error %q should mention \"dev build\"", ver, err.Error())
		}
	}
}

// ----------------------------------------------------------------------------
// TestUpdateResult_Fields
// ----------------------------------------------------------------------------

func TestUpdateResult_Fields(t *testing.T) {
	tests := []struct {
		name        string
		result      updater.UpdateResult
		wantFound   bool
		wantCurrent string
		wantLatest  string
	}{
		{
			name: "update available",
			result: updater.UpdateResult{
				CurrentVersion: "0.1.0",
				LatestVersion:  "0.2.0",
				UpdateFound:    true,
			},
			wantFound:   true,
			wantCurrent: "0.1.0",
			wantLatest:  "0.2.0",
		},
		{
			name: "already on latest",
			result: updater.UpdateResult{
				CurrentVersion: "0.2.0",
				LatestVersion:  "0.2.0",
				UpdateFound:    false,
			},
			wantFound:   false,
			wantCurrent: "0.2.0",
			wantLatest:  "0.2.0",
		},
		{
			name: "no release found",
			result: updater.UpdateResult{
				CurrentVersion: "0.1.0",
				LatestVersion:  "",
				UpdateFound:    false,
			},
			wantFound:   false,
			wantCurrent: "0.1.0",
			wantLatest:  "",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if tc.result.UpdateFound != tc.wantFound {
				t.Errorf("UpdateFound: got %v, want %v", tc.result.UpdateFound, tc.wantFound)
			}
			if tc.result.CurrentVersion != tc.wantCurrent {
				t.Errorf("CurrentVersion: got %q, want %q", tc.result.CurrentVersion, tc.wantCurrent)
			}
			if tc.result.LatestVersion != tc.wantLatest {
				t.Errorf("LatestVersion: got %q, want %q", tc.result.LatestVersion, tc.wantLatest)
			}
		})
	}
}
