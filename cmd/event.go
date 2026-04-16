package cmd

import "github.com/spf13/cobra"

var eventCmd = &cobra.Command{
	Use:   "event",
	Short: "Manage event directives in a ledger",
}

func init() {
	rootCmd.AddCommand(eventCmd)
}
