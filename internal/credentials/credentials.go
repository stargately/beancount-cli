package credentials

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"time"
)

// Credentials is the structure persisted to ~/.beancount/credentials.json.
type Credentials struct {
	Token    string `json:"token"`
	ExpireAt string `json:"expireAt"`
}

func credentialsDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".beancount"), nil
}

func credentialsPath() (string, error) {
	dir, err := credentialsDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "credentials.json"), nil
}

// Save writes token and expiry to ~/.beancount/credentials.json with mode 0600.
func Save(token, expireAt string) error {
	dir, err := credentialsDir()
	if err != nil {
		return err
	}
	if err := os.MkdirAll(dir, 0o700); err != nil {
		return err
	}

	path, err := credentialsPath()
	if err != nil {
		return err
	}

	data, err := json.MarshalIndent(&Credentials{Token: token, ExpireAt: expireAt}, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0o600)
}

// Load reads ~/.beancount/credentials.json. Returns nil, nil when the file does not exist.
func Load() (*Credentials, error) {
	path, err := credentialsPath()
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, nil
		}
		return nil, err
	}

	var creds Credentials
	if err := json.Unmarshal(data, &creds); err != nil {
		return nil, err
	}
	return &creds, nil
}

// Clear deletes ~/.beancount/credentials.json. Silently succeeds if the file is absent.
func Clear() error {
	path, err := credentialsPath()
	if err != nil {
		return err
	}
	if err := os.Remove(path); err != nil && !errors.Is(err, os.ErrNotExist) {
		return err
	}
	return nil
}

// IsExpired reports whether the stored token has passed its expiry timestamp.
func (c *Credentials) IsExpired() bool {
	if c.ExpireAt == "" {
		return true
	}
	t, err := time.Parse(time.RFC3339, c.ExpireAt)
	if err != nil {
		return true
	}
	return time.Now().After(t)
}
