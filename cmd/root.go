package cmd

import (
	"os"

	"beancount.io/beancount-cli/internal/updater"
	"github.com/spf13/cobra"
)

var updateCheckCh <-chan updater.UpdateResult

var rootCmd = &cobra.Command{
	Use:   "beancount-cli",
	Short: "Official CLI for Beancount",
	Long: `beancount-cli is the official command-line interface for Beancount.

Use it to authenticate and interact with your Beancount account.

Environment variables:
  BEANCOUNT_API_URL               Override the GraphQL endpoint (default: https://beancount.io/api-gateway/)
  BEANCOUNT_DASHBOARD_URL         Override the dashboard URL (default: https://beancount.io)
  BEANCOUNT_NO_UPDATE_NOTIFIER    Set to any value to disable update notifications`,
	// Don't print usage on every RunE error — the error message is sufficient.
	SilenceUsage: true,
	PersistentPreRun: func(cmd *cobra.Command, _ []string) {
		if cmd.Name() != "upgrade" {
			updateCheckCh = updater.StartUpdateCheck(cmd.Root().Version)
		}
	},
	PersistentPostRun: func(cmd *cobra.Command, _ []string) {
		if cmd.Name() == "upgrade" || updateCheckCh == nil {
			return
		}
		select {
		case result := <-updateCheckCh:
			updater.PrintUpdateTip(os.Stderr, result)
		default:
		}
	},
}

// SetVersion sets the version string shown by --version.
func SetVersion(v string) {
	rootCmd.Version = v
}

// Execute runs the root command and exits on error.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
