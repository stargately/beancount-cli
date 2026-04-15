package cmd

import "github.com/spf13/cobra"

var transactionCmd = &cobra.Command{
	Use:   "transaction",
	Short: "Manage transactions in a ledger",
}

func init() {
	rootCmd.AddCommand(transactionCmd)
}
