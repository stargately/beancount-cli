package cmd

import (
	"context"
	"fmt"
	"os"

	selfupdate "github.com/creativeprojects/go-selfupdate"
	"github.com/spf13/cobra"
)

var upgradeCheckOnly bool

var upgradeCmd = &cobra.Command{
	Use:   "upgrade",
	Short: "Check for and apply updates to beancount-cli",
	Long: `Check GitHub for the latest beancount-cli release and update the binary in-place.

Use --check to only report the available version without applying the update.`,
	Example: `  beancount-cli upgrade
  beancount-cli upgrade --check`,
	RunE: runUpgrade,
}

func init() {
	rootCmd.AddCommand(upgradeCmd)
	upgradeCmd.Flags().BoolVar(&upgradeCheckOnly, "check", false, "Only check for updates, do not apply")
}

type updateResult struct {
	currentVersion string
	latestVersion  string
	updateFound    bool
	release        *selfupdate.Release
}

// checkUpdate queries GitHub releases for the given slug and compares against
// currentVersion. Returns an error for dev builds or network failures.
func checkUpdate(ctx context.Context, currentVersion, slug string) (updateResult, *selfupdate.Updater, error) {
	if currentVersion == "" || currentVersion == "dev" {
		return updateResult{}, nil, fmt.Errorf("cannot check for updates: version is %q (dev build)", currentVersion)
	}

	updater, err := selfupdate.NewUpdater(selfupdate.Config{})
	if err != nil {
		return updateResult{}, nil, fmt.Errorf("failed to create updater: %w", err)
	}

	latest, found, err := updater.DetectLatest(ctx, selfupdate.ParseSlug(slug))
	if err != nil {
		return updateResult{}, nil, fmt.Errorf("failed to check for updates: %w", err)
	}

	result := updateResult{
		currentVersion: currentVersion,
	}
	if found && latest != nil {
		result.latestVersion = latest.Version()
		result.release = latest
		result.updateFound = latest.GreaterThan(currentVersion)
	}
	return result, updater, nil
}

func runUpgrade(cmd *cobra.Command, _ []string) error {
	result, updater, err := checkUpdate(cmd.Context(), rootCmd.Version, "stargately/beancount-cli")
	if err != nil {
		return err
	}

	fmt.Fprintf(cmd.OutOrStdout(), "Current version: %s\n", result.currentVersion)

	if !result.updateFound {
		fmt.Fprintf(cmd.OutOrStdout(), "Already up to date (%s).\n", result.latestVersion)
		return nil
	}

	fmt.Fprintf(cmd.OutOrStdout(), "Update available: %s → %s\n", result.currentVersion, result.latestVersion)

	if upgradeCheckOnly {
		fmt.Fprintln(cmd.OutOrStdout(), "Run without --check to apply the update.")
		return nil
	}

	exe, err := os.Executable()
	if err != nil {
		return fmt.Errorf("failed to find executable path: %w", err)
	}

	fmt.Fprintln(cmd.OutOrStdout(), "Applying update...")
	if err := updater.UpdateTo(cmd.Context(), result.release, exe); err != nil {
		return fmt.Errorf("update failed: %w", err)
	}
	fmt.Fprintf(cmd.OutOrStdout(), "Successfully updated to %s. Restart beancount-cli to use the new version.\n", result.latestVersion)
	return nil
}
