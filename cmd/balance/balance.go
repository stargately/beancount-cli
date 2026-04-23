package balance

import "github.com/spf13/cobra"

func NewCmdBalance() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "balance",
		Short: "Manage balance directives in a ledger",
	}
	cmd.AddCommand(NewCmdAdd())
	return cmd
}
