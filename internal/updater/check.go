package updater

import (
	"context"
	"fmt"

	selfupdate "github.com/creativeprojects/go-selfupdate"
	goversion "github.com/hashicorp/go-version"
)

// UpdateResult holds the outcome of a version check.
type UpdateResult struct {
	CurrentVersion string
	LatestVersion  string
	UpdateFound    bool
	Release        *selfupdate.Release
}

// CheckUpdate queries GitHub releases for the given slug and compares against currentVersion.
func CheckUpdate(ctx context.Context, currentVersion, slug string) (UpdateResult, *selfupdate.Updater, error) {
	if currentVersion == "" || currentVersion == "dev" {
		return UpdateResult{}, nil, fmt.Errorf("cannot check for updates: version is %q (dev build)", currentVersion)
	}

	su, err := selfupdate.NewUpdater(selfupdate.Config{})
	if err != nil {
		return UpdateResult{}, nil, fmt.Errorf("failed to create updater: %w", err)
	}

	latest, found, err := su.DetectLatest(ctx, selfupdate.ParseSlug(slug))
	if err != nil {
		return UpdateResult{}, nil, fmt.Errorf("failed to check for updates: %w", err)
	}

	result := UpdateResult{CurrentVersion: currentVersion}
	if found && latest != nil {
		result.LatestVersion = latest.Version()
		result.Release = latest
		result.UpdateFound = latest.GreaterThan(currentVersion)
	}
	return result, su, nil
}

// IsNewerVersion reports whether latest is semantically greater than current.
func IsNewerVersion(latest, current string) bool {
	l, err1 := goversion.NewVersion(latest)
	c, err2 := goversion.NewVersion(current)
	if err1 != nil || err2 != nil {
		return false
	}
	return l.GreaterThan(c)
}
