package cmd

import (
	"fmt"
	"os"

	"beancount.io/beancount-cli/internal/updater"
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

func runUpgrade(cmd *cobra.Command, _ []string) error {
	result, su, err := updater.CheckUpdate(cmd.Context(), rootCmd.Version, "stargately/beancount-cli")
	if err != nil {
		return err
	}

	fmt.Fprintf(cmd.OutOrStdout(), "Current version: %s\n", result.CurrentVersion)

	if !result.UpdateFound {
		fmt.Fprintf(cmd.OutOrStdout(), "Already up to date (%s).\n", result.LatestVersion)
		return nil
	}

	fmt.Fprintf(cmd.OutOrStdout(), "Update available: %s → %s\n", result.CurrentVersion, result.LatestVersion)

	if upgradeCheckOnly {
		fmt.Fprintln(cmd.OutOrStdout(), "Run without --check to apply the update.")
		return nil
	}

	exe, err := os.Executable()
	if err != nil {
		return fmt.Errorf("failed to find executable path: %w", err)
	}

	fmt.Fprintln(cmd.OutOrStdout(), "Applying update...")
	if err := su.UpdateTo(cmd.Context(), result.Release, exe); err != nil {
		return fmt.Errorf("update failed: %w", err)
	}
	fmt.Fprintf(cmd.OutOrStdout(), "Successfully updated to %s. Restart beancount-cli to use the new version.\n", result.LatestVersion)
	return nil
}
