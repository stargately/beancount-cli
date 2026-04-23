package commodity

import "github.com/spf13/cobra"

func NewCmdCommodity() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "commodity",
		Short: "Manage commodity directives in a ledger",
	}
	cmd.AddCommand(NewCmdAdd())
	return cmd
}
