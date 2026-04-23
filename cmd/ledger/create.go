package ledger

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"beancount.io/beancount-cli/generated"
	"beancount.io/beancount-cli/internal/utils"
)

func NewCmdCreate() *cobra.Command {
	var name, desc string

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a new ledger",
		Long: `Create a new Beancount ledger.

The ledger name must be in slug format: lowercase letters, digits, hyphens, and
underscores only (e.g. "my-budget", "personal_2024").`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			client, err := utils.NewAuthedClient()
			if err != nil {
				return err
			}
			resp, err := generated.CreateLedger(context.Background(), client, name, &desc)
			if err != nil {
				return fmt.Errorf("failed to create ledger: %w", err)
			}
			ledger := resp.CreateLedger
			fmt.Fprintf(cmd.OutOrStdout(), "Created ledger: %s\n", ledger.Name)
			fmt.Fprintf(cmd.OutOrStdout(), "ID:             %s\n", ledger.Id)
			fmt.Fprintf(cmd.OutOrStdout(), "Full name:      %s\n", ledger.FullName)
			return nil
		},
	}

	cmd.Flags().StringVar(&name, "name", "", "Ledger name in slug format (required)")
	cmd.Flags().StringVar(&desc, "description", "", "Ledger description (optional)")
	_ = cmd.MarkFlagRequired("name")
	return cmd
}
