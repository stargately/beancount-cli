package collaborator

import (
	"context"
	"fmt"
	"strings"
	"text/tabwriter"

	"github.com/spf13/cobra"

	"beancount.io/beancount-cli/generated"
	"beancount.io/beancount-cli/internal/utils"
)

func NewCmdList() *cobra.Command {
	return &cobra.Command{
		Use:   "list <ledger>",
		Short: "List collaborators for a ledger",
		Args:  cobra.ExactArgs(1),
		RunE:  runList,
	}
}

func runList(cmd *cobra.Command, args []string) error {
	client, err := utils.NewAuthedClient()
	if err != nil {
		return err
	}

	ledgerID := utils.LedgerID(args[0])
	resp, err := generated.ListLedgerCollaborators(context.Background(), client, ledgerID)
	if err != nil {
		return fmt.Errorf("failed to list collaborators: %w", err)
	}

	collaborators := resp.ListLedgerCollaborators
	if len(collaborators) == 0 {
		fmt.Fprintln(cmd.OutOrStdout(), "No collaborators found.")
		return nil
	}

	w := tabwriter.NewWriter(cmd.OutOrStdout(), 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "USERNAME\tPERMISSION")
	fmt.Fprintln(w, strings.Repeat("-", 8)+"\t"+strings.Repeat("-", 10))
	for _, c := range collaborators {
		login := ""
		if c.Login != nil {
			login = *c.Login
		}
		permission := ""
		if c.Permission != nil {
			permission = *c.Permission
		}
		fmt.Fprintf(w, "%s\t%s\n", login, permission)
	}
	return w.Flush()
}
