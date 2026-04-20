package updater

import (
	"context"
	"fmt"
	"io"
	"os"
	"time"

	"beancount.io/beancount-cli/internal/config"
)

// StartUpdateCheck launches a background goroutine and returns a channel that
// receives at most one UpdateResult. Uses the 24h state cache to avoid hitting
// the GitHub API on every invocation.
func StartUpdateCheck(version string) <-chan UpdateResult {
	ch := make(chan UpdateResult, 1)
	if !ShouldCheck() || !isStderrTerminal() || version == "" || version == "dev" {
		return ch
	}

	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		state, err := LoadState()
		if err != nil || state == nil {
			state = &State{}
		}

		var result UpdateResult
		if IsStale(state) {
			r, _, err := CheckUpdate(ctx, version, config.RepoSlug)
			if err == nil {
				result = r
				_ = SaveState(&State{
					LastCheckedAt: time.Now(),
					LatestVersion: r.LatestVersion,
				})
			}
		} else if state.LatestVersion != "" {
			result = UpdateResult{
				CurrentVersion: version,
				LatestVersion:  state.LatestVersion,
				UpdateFound:    IsNewerVersion(state.LatestVersion, version),
			}
		}

		ch <- result
	}()

	return ch
}

// PrintUpdateTip writes the upgrade notification to w when a newer version is available.
func PrintUpdateTip(w io.Writer, result UpdateResult) {
	if !result.UpdateFound {
		return
	}
	fmt.Fprintf(w, "\nA new release of beancount-cli is available: %s → %s\n", result.CurrentVersion, result.LatestVersion)
	fmt.Fprintf(w, "To upgrade, run: %s\n", UpgradeCommand())
	fmt.Fprintf(w, "https://github.com/%s/%s/releases/tag/v%s\n\n", config.RepoOwner, config.RepoName, result.LatestVersion)
}

func isStderrTerminal() bool {
	fi, err := os.Stderr.Stat()
	return err == nil && (fi.Mode()&os.ModeCharDevice) != 0
}
