package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "beancount-cli",
	Short: "Official CLI for Beancount",
	Long: `beancount-cli is the official command-line interface for Beancount.

Use it to authenticate and interact with your Beancount account.

Environment variables:
  BEANCOUNT_API_URL        Override the GraphQL endpoint (default: https://beancount.io/api-gateway/)
  BEANCOUNT_DASHBOARD_URL  Override the dashboard URL (default: https://beancount.io)`,
	// Don't print usage on every RunE error — the error message is sufficient.
	SilenceUsage: true,
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
