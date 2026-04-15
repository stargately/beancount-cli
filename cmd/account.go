package cmd

import "github.com/spf13/cobra"

var accountCmd = &cobra.Command{
	Use:   "account",
	Short: "Manage accounts in a ledger",
}

func init() {
	rootCmd.AddCommand(accountCmd)
}
