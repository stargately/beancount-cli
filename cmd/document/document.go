package document

import "github.com/spf13/cobra"

func NewCmdDocument() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "document",
		Short: "Manage document directives in a ledger",
	}
	cmd.AddCommand(NewCmdAdd())
	return cmd
}
