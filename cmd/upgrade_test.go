package cmd

import (
	"context"
	"strings"
	"testing"

	selfupdate "github.com/creativeprojects/go-selfupdate"
)

// ----------------------------------------------------------------------------
// TestCheckUpdate_DevBuild
// ----------------------------------------------------------------------------

func TestCheckUpdate_DevBuild(t *testing.T) {
	for _, ver := range []string{"dev", ""} {
		_, _, err := checkUpdate(context.Background(), ver, "stargately/beancount-cli")
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
		result      updateResult
		wantFound   bool
		wantCurrent string
		wantLatest  string
	}{
		{
			name: "update available",
			result: updateResult{
				currentVersion: "0.1.0",
				latestVersion:  "0.2.0",
				updateFound:    true,
				release:        &selfupdate.Release{},
			},
			wantFound:   true,
			wantCurrent: "0.1.0",
			wantLatest:  "0.2.0",
		},
		{
			name: "already on latest",
			result: updateResult{
				currentVersion: "0.2.0",
				latestVersion:  "0.2.0",
				updateFound:    false,
			},
			wantFound:   false,
			wantCurrent: "0.2.0",
			wantLatest:  "0.2.0",
		},
		{
			name: "no release found",
			result: updateResult{
				currentVersion: "0.1.0",
				latestVersion:  "",
				updateFound:    false,
			},
			wantFound:   false,
			wantCurrent: "0.1.0",
			wantLatest:  "",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if tc.result.updateFound != tc.wantFound {
				t.Errorf("updateFound: got %v, want %v", tc.result.updateFound, tc.wantFound)
			}
			if tc.result.currentVersion != tc.wantCurrent {
				t.Errorf("currentVersion: got %q, want %q", tc.result.currentVersion, tc.wantCurrent)
			}
			if tc.result.latestVersion != tc.wantLatest {
				t.Errorf("latestVersion: got %q, want %q", tc.result.latestVersion, tc.wantLatest)
			}
		})
	}
}
