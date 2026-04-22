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
	transactionListCmd.Flags().IntVar(&txListLimit, "limit", 20, "Maximum number of transactions to fetch")
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
	limit := float64(txListLimit)
	query := &generated.JournalQueryInput{
		DirectiveTypes: []string{"Transaction"},
		Limit:          &limit,
	}
	if txListFilter != "" {
		query.Filter = &txListFilter
	}
	if txListAccount != "" {
		query.Account = &txListAccount
	}

	resp, err := generated.GetLedgerJournal(context.Background(), client, ledgerID, query)
	if err != nil {
		return fmt.Errorf("failed to fetch transactions: %w", err)
	}

	txs, err := tui.DecodeTransactions(resp.GetLedgerJournal.Data)
	if err != nil {
		return fmt.Errorf("failed to decode transactions: %w", err)
	}

	if isTerminal() && !txListNoInteractive {
		return runInteractive(txListLedger, ledgerID, client, query, txs)
	}
	return runPlain(cmd, txs)
}

func runInteractive(ledgerName, ledgerID string, client graphql.Client, query *generated.JournalQueryInput, txs []tui.Transaction) error {
	m := tui.NewTransactionList(ledgerName, ledgerID, client, query, txs)
	p := tea.NewProgram(m, tea.WithAltScreen())
	_, err := p.Run()
	return err
}

func isTerminal() bool {
	fi, err := os.Stdout.Stat()
	if err != nil {
		return false
	}
	return (fi.Mode() & os.ModeCharDevice) != 0
}

func runPlain(cmd *cobra.Command, txs []tui.Transaction) error {
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
