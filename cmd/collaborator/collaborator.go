package collaborator

import "github.com/spf13/cobra"

func NewCmdCollaborator() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "collaborator",
		Short: "Manage ledger collaborators",
	}
	cmd.AddCommand(NewCmdList())
	cmd.AddCommand(NewCmdAdd())
	cmd.AddCommand(NewCmdRemove())
	return cmd
}
