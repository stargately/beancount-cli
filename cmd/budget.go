package cmd

import "github.com/spf13/cobra"

var budgetCmd = &cobra.Command{
	Use:   "budget",
	Short: "Manage budget directives in a ledger",
}

func init() {
	rootCmd.AddCommand(budgetCmd)
}
