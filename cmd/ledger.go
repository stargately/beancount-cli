package cmd

import "github.com/spf13/cobra"

var ledgerCmd = &cobra.Command{
	Use:   "ledger",
	Short: "Manage your Beancount ledgers",
}

func init() {
	rootCmd.AddCommand(ledgerCmd)
}
