package cmd

import "github.com/spf13/cobra"

var balanceCmd = &cobra.Command{
	Use:   "balance",
	Short: "Manage balance directives in a ledger",
}

func init() {
	rootCmd.AddCommand(balanceCmd)
}
