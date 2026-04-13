package cmd

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"beancount.io/beancount-cli/generated"
	"beancount.io/beancount-cli/internal/config"
	"beancount.io/beancount-cli/internal/credentials"
	"beancount.io/beancount-cli/internal/gqlclient"
)

var ledgerDeleteCmd = &cobra.Command{
	Use:   "delete <ledger-id>",
	Short: "Delete a ledger",
	Long: `Delete a ledger by its ID.

The ledger ID can be found with 'beancount-cli ledger list'.

WARNING: This action is irreversible. All ledger data will be permanently deleted.`,
	Args: cobra.ExactArgs(1),
	RunE: runLedgerDelete,
}

func init() {
	ledgerCmd.AddCommand(ledgerDeleteCmd)
}

func runLedgerDelete(cmd *cobra.Command, args []string) error {
	creds, err := credentials.Load()
	if err != nil {
		return fmt.Errorf("failed to read credentials: %w", err)
	}
	if creds == nil || creds.Token == "" {
		return fmt.Errorf("not logged in. Run 'beancount-cli login' to authenticate")
	}
	if creds.IsExpired() {
		return fmt.Errorf("your session has expired. Run 'beancount-cli login' to re-authenticate")
	}

	cfg := config.Load()
	client := gqlclient.NewAuthed(cfg.APIURL, creds.Token)

	ledgerID := args[0]
	_, err = generated.DeleteLedger(context.Background(), client, ledgerID)
	if err != nil {
		return fmt.Errorf("failed to delete ledger: %w", err)
	}

	fmt.Fprintf(cmd.OutOrStdout(), "Deleted ledger: %s\n", ledgerID)
	return nil
}
