package cmd

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"

	"github.com/Khan/genqlient/graphql"

	"beancount.io/beancount-cli/generated"
	"beancount.io/beancount-cli/internal/tools"
	"beancount.io/beancount-cli/internal/utils"
)

// beanQueryCmd — local bean-query tool (unchanged behavior)

var beanQueryFormat string

var beanQueryCmd = &cobra.Command{
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
	RunE: runBeanQuery,
}

func init() {
	rootCmd.AddCommand(beanQueryCmd)
	beanQueryCmd.Flags().StringVarP(&beanQueryFormat, "format", "f", "text", `output format: "text" or "csv"`)

	rootCmd.AddCommand(queryCmd)
}

func runBeanQuery(cmd *cobra.Command, args []string) error {
	if beanQueryFormat != "text" && beanQueryFormat != "csv" {
		return fmt.Errorf("invalid format %q: must be \"text\" or \"csv\"", beanQueryFormat)
	}

	filename := args[0]
	bqlQuery := ""
	if len(args) == 2 {
		bqlQuery = args[1]
	}

	opts := tools.QueryOptions{Format: beanQueryFormat}
	err := tools.Query(context.Background(), filename, bqlQuery, opts)
	if err == nil {
		return nil
	}

	if e, ok := err.(*exec.ExitError); ok {
		os.Exit(e.ExitCode())
	}

	return err
}

// queryCmd — remote BQL query against a hosted ledger

var queryCmd = &cobra.Command{
	Use:   "query <ledger> [bql]",
	Short: "Run a BQL query against a remote ledger",
	Long: `Run a BQL query against a remote Beancount ledger.

With a BQL argument, runs in batch mode and prints the result:
  beancount-cli query user/mybook "SELECT account, sum(position) GROUP BY account"

Without a BQL argument, opens an interactive prompt:
  beancount-cli query user/mybook`,
	Args: cobra.RangeArgs(1, 2),
	RunE: runQuery,
}

func runQuery(cmd *cobra.Command, args []string) error {
	ledger := args[0]
	ledgerID := utils.LedgerID(ledger)

	client, err := utils.NewAuthedClient()
	if err != nil {
		return err
	}

	if len(args) == 2 {
		return execQuery(cmd, client, ledgerID, args[1])
	}

	// Interactive REPL
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Fprint(cmd.OutOrStdout(), "bql> ")
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "exit" || line == "quit" {
			break
		}
		if line != "" {
			if err := execQuery(cmd, client, ledgerID, line); err != nil {
				fmt.Fprintf(cmd.ErrOrStderr(), "error: %v\n", err)
			}
		}
		fmt.Fprint(cmd.OutOrStdout(), "bql> ")
	}
	return scanner.Err()
}

func execQuery(cmd *cobra.Command, client graphql.Client, ledgerID, bql string) error {
	resp, err := generated.QueryShellText(context.Background(), client, ledgerID, bql)
	if err != nil {
		return fmt.Errorf("query failed: %w", err)
	}
	if resp.QueryShellText == nil {
		return fmt.Errorf("no result returned")
	}
	fmt.Fprint(cmd.OutOrStdout(), resp.QueryShellText.Text)
	if !strings.HasSuffix(resp.QueryShellText.Text, "\n") {
		fmt.Fprintln(cmd.OutOrStdout())
	}
	return nil
}
