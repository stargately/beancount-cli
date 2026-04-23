package ledger

import "github.com/spf13/cobra"

func NewCmdLedger() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ledger",
		Short: "Manage your Beancount ledgers",
	}
	cmd.AddCommand(NewCmdCreate())
	cmd.AddCommand(NewCmdList())
	cmd.AddCommand(NewCmdDelete())
	return cmd
}
