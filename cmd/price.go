package cmd

import "github.com/spf13/cobra"

var priceCmd = &cobra.Command{
	Use:   "price",
	Short: "Manage price directives in a ledger",
}

func init() {
	rootCmd.AddCommand(priceCmd)
}
