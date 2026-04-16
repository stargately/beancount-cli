package cmd

import "github.com/spf13/cobra"

var noteCmd = &cobra.Command{
	Use:   "note",
	Short: "Manage note directives in a ledger",
}

func init() {
	rootCmd.AddCommand(noteCmd)
}
