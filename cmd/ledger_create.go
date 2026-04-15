package cmd

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"beancount.io/beancount-cli/generated"
	"beancount.io/beancount-cli/internal/utils"
)

var ledgerCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new ledger",
	Long: `Create a new Beancount ledger.

The ledger name must be in slug format: lowercase letters, digits, hyphens, and
underscores only (e.g. "my-budget", "personal_2024").`,
	RunE: runLedgerCreate,
}

var (
	ledgerCreateName string
	ledgerCreateDesc string
)

func init() {
	ledgerCmd.AddCommand(ledgerCreateCmd)
	ledgerCreateCmd.Flags().StringVar(&ledgerCreateName, "name", "", "Ledger name in slug format (required)")
	ledgerCreateCmd.Flags().StringVar(&ledgerCreateDesc, "description", "", "Ledger description (optional)")
	_ = ledgerCreateCmd.MarkFlagRequired("name")
}

func runLedgerCreate(cmd *cobra.Command, _ []string) error {
	client, err := utils.NewAuthedClient()
	if err != nil {
		return err
	}

	resp, err := generated.CreateLedger(context.Background(), client, ledgerCreateName, &ledgerCreateDesc)
	if err != nil {
		return fmt.Errorf("failed to create ledger: %w", err)
	}

	ledger := resp.CreateLedger
	fmt.Fprintf(cmd.OutOrStdout(), "Created ledger: %s\n", ledger.Name)
	fmt.Fprintf(cmd.OutOrStdout(), "ID:             %s\n", ledger.Id)
	fmt.Fprintf(cmd.OutOrStdout(), "Full name:      %s\n", ledger.FullName)
	return nil
}
