package list

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
	"beancount.io/beancount-cli/internal/utils"
)

func NewCmdList() *cobra.Command {
	var ledger string
	var limit int
	var filter, account string
	var noInteractive bool

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List transactions in a ledger",
		RunE: func(cmd *cobra.Command, _ []string) error {
			client, err := utils.NewAuthedClient()
			if err != nil {
				return err
			}

			ledgerID := utils.LedgerID(ledger)

			baseQuery := &generated.JournalQueryInput{
				DirectiveTypes: []string{"Transaction"},
			}
			if filter != "" {
				baseQuery.Filter = &filter
			}
			if account != "" {
				baseQuery.Account = &account
			}

			if isTerminal() && !noInteractive {
				return runInteractive(cmd, client, ledgerID, ledger, baseQuery)
			}
			return runPlain(cmd, client, ledgerID, baseQuery, limit)
		},
	}

	cmd.Flags().StringVarP(&ledger, "ledger", "l", "", "Ledger full name (e.g. user/mybook) (required)")
	cmd.Flags().IntVar(&limit, "limit", 100, "Maximum number of transactions to fetch (plain mode only)")
	cmd.Flags().StringVar(&filter, "filter", "", "Beancount filter expression (optional)")
	cmd.Flags().StringVar(&account, "account", "", "Filter by account (optional)")
	cmd.Flags().BoolVar(&noInteractive, "no-interactive", false, "Force plain table output instead of TUI")
	_ = cmd.MarkFlagRequired("ledger")
	return cmd
}

func runInteractive(_ *cobra.Command, client graphql.Client, ledgerID, ledgerName string, baseQuery *generated.JournalQueryInput) error {
	initLimit := float64(PageSize)
	q := *baseQuery
	q.Limit = &initLimit

	resp, err := generated.GetLedgerJournal(context.Background(), client, ledgerID, &q)
	if err != nil {
		return fmt.Errorf("failed to fetch transactions: %w", err)
	}
	txs, err := DecodeTransactions(resp.GetLedgerJournal.Data)
	if err != nil {
		return fmt.Errorf("failed to decode transactions: %w", err)
	}

	m := NewTransactionList(ledgerName, ledgerID, client, baseQuery, txs, int(resp.GetLedgerJournal.Total))
	p := tea.NewProgram(m, tea.WithAltScreen())
	_, err = p.Run()
	return err
}

func runPlain(cmd *cobra.Command, client graphql.Client, ledgerID string, baseQuery *generated.JournalQueryInput, limit int) error {
	limitF := float64(limit)
	q := *baseQuery
	q.Limit = &limitF

	resp, err := generated.GetLedgerJournal(context.Background(), client, ledgerID, &q)
	if err != nil {
		return fmt.Errorf("failed to fetch transactions: %w", err)
	}
	txs, err := DecodeTransactions(resp.GetLedgerJournal.Data)
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
