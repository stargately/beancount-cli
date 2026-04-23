package collaborator

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"beancount.io/beancount-cli/generated"
	"beancount.io/beancount-cli/internal/utils"
)

func NewCmdAdd() *cobra.Command {
	var permission string

	cmd := &cobra.Command{
		Use:   "add <ledger> <username>",
		Short: "Add or update a collaborator on a ledger",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			ledger, username := args[0], args[1]
			ledgerID, collaborator, perm := buildAddInput(ledger, username, permission)

			client, err := utils.NewAuthedClient()
			if err != nil {
				return err
			}

			resp, err := generated.AddOrUpdateLedgerCollaborator(context.Background(), client, ledgerID, collaborator, &perm)
			if err != nil {
				return fmt.Errorf("failed to add collaborator: %w", err)
			}

			result := resp.AddOrUpdateLedgerCollaborator
			if !result.Success {
				msg := "unknown error"
				if result.Message != nil && *result.Message != "" {
					msg = *result.Message
				}
				return fmt.Errorf("server error: %s", msg)
			}

			fmt.Fprintf(cmd.OutOrStdout(), "Added %s to %s with %s permission\n", username, ledger, perm)
			return nil
		},
	}

	cmd.Flags().StringVar(&permission, "permission", "read", "Permission level (read or write)")
	return cmd
}

func buildAddInput(ledger, username, permission string) (ledgerID, collaborator, perm string) {
	if permission == "" {
		permission = "read"
	}
	return utils.LedgerID(ledger), username, permission
}
