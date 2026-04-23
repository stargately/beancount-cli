package ledger

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
		Use:   "list",
		Short: "List your ledgers",
		RunE:  runList,
	}
}

func runList(cmd *cobra.Command, _ []string) error {
	client, err := utils.NewAuthedClient()
	if err != nil {
		return err
	}

	resp, err := generated.ListUserOwnedLedgers(context.Background(), client)
	if err != nil {
		return fmt.Errorf("failed to list ledgers: %w", err)
	}

	ledgers := resp.ListUserOwnedLedgers
	if len(ledgers) == 0 {
		fmt.Fprintln(cmd.OutOrStdout(), "No ledgers found.")
		return nil
	}

	w := tabwriter.NewWriter(cmd.OutOrStdout(), 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "NAME\tFULLNAME\tDESCRIPTION\tPRIVATE\tCREATED")
	fmt.Fprintln(w, strings.Repeat("-", 4)+"\t"+strings.Repeat("-", 8)+"\t"+strings.Repeat("-", 11)+"\t"+strings.Repeat("-", 7)+"\t"+strings.Repeat("-", 7))
	for _, l := range ledgers {
		desc := "(none)"
		if l.Description != nil && *l.Description != "" {
			desc = *l.Description
		}
		private := "no"
		if l.Private {
			private = "yes"
		}
		created := l.CreatedAt
		if len(created) >= 10 {
			created = created[:10]
		}
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\n", l.Name, l.FullName, desc, private, created)
	}
	return w.Flush()
}
