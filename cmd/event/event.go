package event

import "github.com/spf13/cobra"

func NewCmdEvent() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "event",
		Short: "Manage event directives in a ledger",
	}
	cmd.AddCommand(NewCmdAdd())
	return cmd
}
