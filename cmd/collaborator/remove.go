package collaborator

import (
	"bufio"
	"context"
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"beancount.io/beancount-cli/generated"
	"beancount.io/beancount-cli/internal/utils"
)

func NewCmdRemove() *cobra.Command {
	return &cobra.Command{
		Use:   "remove <ledger> <username>",
		Short: "Remove a collaborator from a ledger",
		Args:  cobra.ExactArgs(2),
		RunE:  runRemove,
	}
}

func runRemove(cmd *cobra.Command, args []string) error {
	ledger, username := args[0], args[1]

	fmt.Fprintf(cmd.OutOrStdout(), "Are you sure you want to remove %s from %s? [y/N]: ", username, ledger)
	reader := bufio.NewReader(cmd.InOrStdin())
	answer, err := reader.ReadString('\n')
	if err != nil {
		return fmt.Errorf("failed to read input: %w", err)
	}
	answer = strings.TrimSpace(strings.ToLower(answer))
	if answer != "y" && answer != "yes" {
		fmt.Fprintln(cmd.OutOrStdout(), "Aborted.")
		return nil
	}

	client, err := utils.NewAuthedClient()
	if err != nil {
		return err
	}

	ledgerID := utils.LedgerID(ledger)
	resp, err := generated.DeleteLedgerCollaborator(context.Background(), client, ledgerID, username)
	if err != nil {
		return fmt.Errorf("failed to remove collaborator: %w", err)
	}

	result := resp.DeleteLedgerCollaborator
	if !result.Success {
		msg := "unknown error"
		if result.Message != nil && *result.Message != "" {
			msg = *result.Message
		}
		return fmt.Errorf("server error: %s", msg)
	}

	fmt.Fprintf(cmd.OutOrStdout(), "Removed %s from %s\n", username, ledger)
	return nil
}
