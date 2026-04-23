package budget

import "github.com/spf13/cobra"

func NewCmdBudget() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "budget",
		Short: "Manage budget directives in a ledger",
	}
	cmd.AddCommand(NewCmdAdd())
	return cmd
}
