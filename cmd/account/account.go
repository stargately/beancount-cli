package account

import "github.com/spf13/cobra"

func NewCmdAccount() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "account",
		Short: "Manage accounts in a ledger",
	}
	cmd.AddCommand(NewCmdOpen())
	cmd.AddCommand(NewCmdClose())
	return cmd
}
