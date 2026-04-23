package note

import "github.com/spf13/cobra"

func NewCmdNote() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "note",
		Short: "Manage note directives in a ledger",
	}
	cmd.AddCommand(NewCmdAdd())
	return cmd
}
