package cmd

import (
	"context"
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"

	"beancount.io/beancount-cli/internal/tools"
)

var queryFormat string

var queryCmd = &cobra.Command{
	Use:   "bean-query <file> [query]",
	Short: "Query a Beancount ledger using BQL",
	Long: `bean-query runs bean-query against a Beancount ledger file.

With a query argument, bean-query runs in batch mode and prints results to stdout:
  beancount-cli bean-query ledger.beancount "SELECT account, sum(position)"

Without a query argument, bean-query opens its interactive shell where you can
type BQL statements at the prompt:
  beancount-cli bean-query ledger.beancount

Output format can be changed with --format:
  beancount-cli bean-query --format csv ledger.beancount "BALANCES"`,
	Args: cobra.RangeArgs(1, 2),
	RunE: runQuery,
}

func init() {
	rootCmd.AddCommand(queryCmd)
	queryCmd.Flags().StringVarP(&queryFormat, "format", "f", "text", `output format: "text" or "csv"`)
}

func runQuery(cmd *cobra.Command, args []string) error {
	if queryFormat != "text" && queryFormat != "csv" {
		return fmt.Errorf("invalid format %q: must be \"text\" or \"csv\"", queryFormat)
	}

	filename := args[0]
	bqlQuery := ""
	if len(args) == 2 {
		bqlQuery = args[1]
	}

	opts := tools.QueryOptions{Format: queryFormat}
	err := tools.Query(context.Background(), filename, bqlQuery, opts)
	if err == nil {
		return nil
	}

	// bean-query already wrote its output/errors to stdout/stderr.
	// Exit with the same code rather than printing a redundant error message.
	if e, ok := err.(*exec.ExitError); ok {
		os.Exit(e.ExitCode())
	}

	// Tool not found, file not readable, etc. — surface the message.
	return err
}
