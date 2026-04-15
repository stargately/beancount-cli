package cmd

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"beancount.io/beancount-cli/generated"
	"beancount.io/beancount-cli/internal/utils"
)

var ledgerDeleteCmd = &cobra.Command{
	Use:   "delete <fullname>",
	Short: "Delete a ledger",
	Long: `Delete a ledger by its full name.

The full name can be found with 'beancount-cli ledger list'.

WARNING: This action is irreversible. All ledger data will be permanently deleted.`,
	Args: cobra.ExactArgs(1),
	RunE: runLedgerDelete,
}

func init() {
	ledgerCmd.AddCommand(ledgerDeleteCmd)
}

func runLedgerDelete(cmd *cobra.Command, args []string) error {
	client, err := utils.NewAuthedClient()
	if err != nil {
		return err
	}

	fullName := args[0]
	resp, err := generated.ListUserOwnedLedgers(context.Background(), client)
	if err != nil {
		return fmt.Errorf("failed to list ledgers: %w", err)
	}

	var ledgerID string
	for _, l := range resp.ListUserOwnedLedgers {
		if l.FullName == fullName {
			ledgerID = l.Id
			break
		}
	}
	if ledgerID == "" {
		return fmt.Errorf("ledger not found: %s", fullName)
	}

	_, err = generated.DeleteLedger(context.Background(), client, ledgerID)
	if err != nil {
		return fmt.Errorf("failed to delete ledger: %w", err)
	}

	fmt.Fprintf(cmd.OutOrStdout(), "Deleted ledger: %s\n", fullName)
	return nil
}
