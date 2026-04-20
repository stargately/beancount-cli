package cmd

import "github.com/spf13/cobra"

var commodityCmd = &cobra.Command{
	Use:   "commodity",
	Short: "Manage commodity directives in a ledger",
}

func init() {
	rootCmd.AddCommand(commodityCmd)
}
