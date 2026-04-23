package price

import "github.com/spf13/cobra"

func NewCmdPrice() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "price",
		Short: "Manage price directives in a ledger",
	}
	cmd.AddCommand(NewCmdAdd())
	return cmd
}
