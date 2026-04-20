package cmd

import "github.com/spf13/cobra"

var documentCmd = &cobra.Command{
	Use:   "document",
	Short: "Manage document directives in a ledger",
}

func init() {
	rootCmd.AddCommand(documentCmd)
}
