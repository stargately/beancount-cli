package cmd

import (
	"context"
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/Khan/genqlient/graphql"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"

	"beancount.io/beancount-cli/generated"
	"beancount.io/beancount-cli/internal/tui"
	"beancount.io/beancount-cli/internal/utils"
)

var transactionListCmd = &cobra.Command{
	Use:   "list",
	Short: "List transactions in a ledger",
	RunE:  runTransactionList,
}

var (
	txListLedger        string
	txListLimit         int
	txListFilter        string
	txListAccount       string
	txListNoInteractive bool
)

func init() {
	transactionCmd.AddCommand(transactionListCmd)
	transactionListCmd.Flags().StringVarP(&txListLedger, "ledger", "l", "", "Ledger full name (e.g. user/mybook) (required)")
	transactionListCmd.Flags().IntVar(&txListLimit, "limit", 100, "Maximum number of transactions to fetch (plain mode only)")
	transactionListCmd.Flags().StringVar(&txListFilter, "filter", "", "Beancount filter expression (optional)")
	transactionListCmd.Flags().StringVar(&txListAccount, "account", "", "Filter by account (optional)")
	transactionListCmd.Flags().BoolVar(&txListNoInteractive, "no-interactive", false, "Force plain table output instead of TUI")
	_ = transactionListCmd.MarkFlagRequired("ledger")
}

func runTransactionList(cmd *cobra.Command, _ []string) error {
	client, err := utils.NewAuthedClient()
	if err != nil {
		return err
	}

	ledgerID := utils.LedgerID(txListLedger)

	// Base query template — no limit/offset; those are set per-call.
	baseQuery := &generated.JournalQueryInput{
		DirectiveTypes: []string{"Transaction"},
	}
	if txListFilter != "" {
		baseQuery.Filter = &txListFilter
	}
	if txListAccount != "" {
		baseQuery.Account = &txListAccount
	}

	if isTerminal() && !txListNoInteractive {
		return runInteractive(cmd, client, ledgerID, baseQuery)
	}
	return runPlain(cmd, client, ledgerID, baseQuery)
}

func runInteractive(_ *cobra.Command, client graphql.Client, ledgerID string, baseQuery *generated.JournalQueryInput) error {
	initLimit := float64(tui.PageSize)
	q := *baseQuery
	q.Limit = &initLimit

	resp, err := generated.GetLedgerJournal(context.Background(), client, ledgerID, &q)
	if err != nil {
		return fmt.Errorf("failed to fetch transactions: %w", err)
	}
	txs, err := tui.DecodeTransactions(resp.GetLedgerJournal.Data)
	if err != nil {
		return fmt.Errorf("failed to decode transactions: %w", err)
	}

	m := tui.NewTransactionList(txListLedger, ledgerID, client, baseQuery, txs, int(resp.GetLedgerJournal.Total))
	p := tea.NewProgram(m, tea.WithAltScreen())
	_, err = p.Run()
	return err
}

func runPlain(cmd *cobra.Command, client graphql.Client, ledgerID string, baseQuery *generated.JournalQueryInput) error {
	limit := float64(txListLimit)
	q := *baseQuery
	q.Limit = &limit

	resp, err := generated.GetLedgerJournal(context.Background(), client, ledgerID, &q)
	if err != nil {
		return fmt.Errorf("failed to fetch transactions: %w", err)
	}
	txs, err := tui.DecodeTransactions(resp.GetLedgerJournal.Data)
	if err != nil {
		return fmt.Errorf("failed to decode transactions: %w", err)
	}

	if len(txs) == 0 {
		fmt.Fprintln(cmd.OutOrStdout(), "No transactions found.")
		return nil
	}

	w := tabwriter.NewWriter(cmd.OutOrStdout(), 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "DATE\tPAYEE\tNARRATION")
	fmt.Fprintln(w, strings.Repeat("-", 10)+"\t"+strings.Repeat("-", 20)+"\t"+strings.Repeat("-", 20))
	for _, tx := range txs {
		fmt.Fprintf(w, "%s\t%s\t%s\n", tx.Date, tx.Payee, tx.Narration)
	}
	return w.Flush()
}

func isTerminal() bool {
	fi, err := os.Stdout.Stat()
	if err != nil {
		return false
	}
	return (fi.Mode() & os.ModeCharDevice) != 0
}
