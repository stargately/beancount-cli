package updater

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const checkTTL = 24 * time.Hour

// State is persisted to ~/.beancount/update_state.json.
type State struct {
	LastCheckedAt time.Time `json:"last_checked_at"`
	LatestVersion string    `json:"latest_version"`
}

func stateDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".beancount"), nil
}

func statePath() (string, error) {
	dir, err := stateDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "update_state.json"), nil
}

// LoadState reads the update state from disk. Returns an empty state on missing file.
func LoadState() (*State, error) {
	path, err := statePath()
	if err != nil {
		return nil, err
	}
	data, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return &State{}, nil
		}
		return nil, err
	}
	var s State
	if err := json.Unmarshal(data, &s); err != nil {
		return &State{}, nil
	}
	return &s, nil
}

// SaveState writes the update state to disk.
func SaveState(s *State) error {
	dir, err := stateDir()
	if err != nil {
		return err
	}
	if err := os.MkdirAll(dir, 0o700); err != nil {
		return err
	}
	path, err := statePath()
	if err != nil {
		return err
	}
	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0o600)
}

// IsStale reports whether the cached state is older than the 24h TTL.
func IsStale(s *State) bool {
	return time.Since(s.LastCheckedAt) > checkTTL
}

// ShouldCheck reports whether the passive update check should run.
// Returns false in CI environments or when BEANCOUNT_NO_UPDATE_NOTIFIER is set.
func ShouldCheck() bool {
	if os.Getenv("BEANCOUNT_NO_UPDATE_NOTIFIER") != "" {
		return false
	}
	if os.Getenv("CI") != "" || os.Getenv("GITHUB_ACTIONS") != "" {
		return false
	}
	return true
}

// IsBrewInstall reports whether the running binary was installed via Homebrew.
func IsBrewInstall() bool {
	exe, err := os.Executable()
	if err != nil {
		return false
	}
	p := strings.ToLower(exe)
	return strings.Contains(p, "/homebrew/") || strings.Contains(p, "/linuxbrew/")
}

// UpgradeCommand returns the appropriate shell command to upgrade beancount-cli.
func UpgradeCommand() string {
	if IsBrewInstall() {
		return "brew upgrade beancount-cli"
	}
	return "beancount-cli upgrade"
}
