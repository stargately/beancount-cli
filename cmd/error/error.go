package cmderror

import (
	"context"
	"fmt"
	"text/tabwriter"

	"github.com/spf13/cobra"

	"beancount.io/beancount-cli/generated"
	"beancount.io/beancount-cli/internal/utils"
)

func NewCmdError() *cobra.Command {
	return &cobra.Command{
		Use:   "error <ledger>",
		Short: "Show errors for a ledger",
		Args:  cobra.ExactArgs(1),
		RunE:  runError,
	}
}

func runError(cmd *cobra.Command, args []string) error {
	client, err := utils.NewAuthedClient()
	if err != nil {
		return err
	}

	ledgerID := utils.LedgerID(args[0])
	resp, err := generated.GetLedgerErrors(context.Background(), client, ledgerID)
	if err != nil {
		return fmt.Errorf("failed to get ledger errors: %w", err)
	}

	errs := resp.GetLedgerErrors
	if len(errs) == 0 {
		fmt.Fprintln(cmd.OutOrStdout(), "No errors found.")
		return nil
	}

	w := tabwriter.NewWriter(cmd.OutOrStdout(), 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "FILE\tLINE\tMESSAGE")
	for _, e := range errs {
		filename := ""
		if e.Filename != nil {
			filename = *e.Filename
		}
		lineno := ""
		if e.Lineno != nil {
			lineno = fmt.Sprintf("%.0f", *e.Lineno)
		}
		fmt.Fprintf(w, "%s\t%s\t%s\n", filename, lineno, e.Message)
	}
	return w.Flush()
}
