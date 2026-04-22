package tui

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"beancount.io/beancount-cli/generated"
	"github.com/Khan/genqlient/graphql"
)

type viewState int

const (
	stateList viewState = iota
	stateDetail
	stateConfirmDelete
	stateDeleting
)

type Transaction struct {
	EntryHash string    `json:"entry_hash"`
	Date      string    `json:"date"`
	Payee     string    `json:"payee"`
	Narration string    `json:"narration"`
	Postings  []Posting `json:"postings"`
	Flag      string    `json:"flag"`
	Tags      []string  `json:"tags"`
	Links     []string  `json:"links"`
}

type Posting struct {
	Account string       `json:"account"`
	Units   PostingUnits `json:"units"`
}

type PostingUnits struct {
	Number   string `json:"number"`
	Currency string `json:"currency"`
}

type contextFetchedMsg struct{ sha256sum string }
type deletedMsg struct{}
type refetchedMsg struct{ items []Transaction }
type errMsg struct{ err error }

type TransactionListModel struct {
	items         []Transaction
	cursor        int
	offset        int
	state         viewState
	ledgerName    string
	ledgerID      string
	client        graphql.Client
	query         *generated.JournalQueryInput
	width         int
	height        int
	loading       bool
	err           error
	status        string
	pendingStatus string
	ctxSha        string
	pendingHash   string
}

func NewTransactionList(ledgerName, ledgerID string, client graphql.Client, query *generated.JournalQueryInput, items []Transaction) TransactionListModel {
	return TransactionListModel{
		items:      items,
		ledgerName: ledgerName,
		ledgerID:   ledgerID,
		client:     client,
		query:      query,
		width:      80,
		height:     24,
	}
}

func DecodeTransactions(raw []any) ([]Transaction, error) {
	b, err := json.Marshal(raw)
	if err != nil {
		return nil, fmt.Errorf("marshal: %w", err)
	}
	var txs []Transaction
	if err := json.Unmarshal(b, &txs); err != nil {
		return nil, fmt.Errorf("decode transactions: %w", err)
	}
	return txs, nil
}

func (m TransactionListModel) Init() tea.Cmd { return nil }

func (m TransactionListModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

	case tea.KeyMsg:
		return m.handleKey(msg)

	case contextFetchedMsg:
		m.ctxSha = msg.sha256sum
		m.state = stateConfirmDelete

	case deletedMsg:
		m.state = stateList
		m.loading = true
		m.pendingStatus = "Transaction deleted."
		m.ctxSha = ""
		m.pendingHash = ""
		return m, m.refetchCmd()

	case refetchedMsg:
		m.items = msg.items
		m.loading = false
		m.status = m.pendingStatus
		m.pendingStatus = ""
		if m.cursor >= len(m.items) && m.cursor > 0 {
			m.cursor--
		}

	case errMsg:
		m.err = msg.err
		m.state = stateDetail
		m.ctxSha = ""
		m.pendingHash = ""
	}

	return m, nil
}

func (m TransactionListModel) handleKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch m.state {
	case stateList:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
				if m.cursor < m.offset {
					m.offset--
				}
			}
		case "down", "j":
			if m.cursor < len(m.items)-1 {
				m.cursor++
				vis := m.listVisibleHeight()
				if m.cursor >= m.offset+vis {
					m.offset++
				}
			}
		case "enter":
			if len(m.items) > 0 {
				m.state = stateDetail
				m.status = ""
				m.err = nil
			}
		case "r":
			if !m.loading {
				m.loading = true
				m.status = ""
				m.err = nil
				return m, m.refetchCmd()
			}
		}

	case stateDetail:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "esc":
			m.state = stateList
			m.err = nil
		case "d":
			if len(m.items) > 0 {
				m.pendingHash = m.items[m.cursor].EntryHash
				return m, m.fetchContextCmd()
			}
		}

	case stateConfirmDelete:
		switch msg.String() {
		case "y":
			m.state = stateDeleting
			return m, m.doDeleteCmd()
		case "n", "esc":
			m.state = stateDetail
			m.ctxSha = ""
			m.pendingHash = ""
		case "q", "ctrl+c":
			return m, tea.Quit
		}
	}

	return m, nil
}

func (m TransactionListModel) fetchContextCmd() tea.Cmd {
	client, ledgerID, entryHash := m.client, m.ledgerID, m.pendingHash
	return func() tea.Msg {
		resp, err := generated.GetLedgerEntryContext(context.Background(), client, ledgerID, entryHash)
		if err != nil {
			return errMsg{err}
		}
		return contextFetchedMsg{resp.GetLedgerEntryContext.Sha256sum}
	}
}

func (m TransactionListModel) doDeleteCmd() tea.Cmd {
	client, ledgerID, entryHash, sha := m.client, m.ledgerID, m.pendingHash, m.ctxSha
	return func() tea.Msg {
		_, err := generated.DeleteLedgerEntrySourceSlice(
			context.Background(), client, ledgerID,
			generated.DeleteSourceSliceInput{EntryHash: entryHash, Sha256sum: sha},
		)
		if err != nil {
			return errMsg{err}
		}
		return deletedMsg{}
	}
}

func (m TransactionListModel) refetchCmd() tea.Cmd {
	client, ledgerID, query := m.client, m.ledgerID, m.query
	return func() tea.Msg {
		resp, err := generated.GetLedgerJournal(context.Background(), client, ledgerID, query)
		if err != nil {
			return errMsg{err}
		}
		items, err := DecodeTransactions(resp.GetLedgerJournal.Data)
		if err != nil {
			return errMsg{err}
		}
		return refetchedMsg{items}
	}
}

func (m TransactionListModel) listVisibleHeight() int {
	return max(m.height-5, 1)
}

var (
	cursorStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("212")).Bold(true)
	selectedStyle = lipgloss.NewStyle().Bold(true)
	dimStyle      = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	boldStyle     = lipgloss.NewStyle().Bold(true)
	errorStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("196"))
	warnStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("214"))
	successStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("82"))
)

func (m TransactionListModel) View() string {
	switch m.state {
	case stateList:
		return m.listView()
	case stateDetail:
		return m.detailView()
	case stateConfirmDelete:
		return m.confirmView()
	case stateDeleting:
		return m.deletingView()
	}
	return ""
}

func (m TransactionListModel) listView() string {
	var sb strings.Builder

	fmt.Fprintf(&sb, "%s\n\n", boldStyle.Render("Transactions: "+m.ledgerName))

	if len(m.items) == 0 {
		sb.WriteString(dimStyle.Render("  No transactions found.") + "\n")
	} else {
		vis := m.listVisibleHeight()
		end := min(m.offset+vis, len(m.items))

		colW := m.width - 4
		dateW := 12
		payeeW := 20
		narrW := max(colW-dateW-payeeW, 10)

		for i := m.offset; i < end; i++ {
			tx := m.items[i]
			row := fmt.Sprintf("%-*s  %-*s  %s", dateW, tx.Date, payeeW, truncate(tx.Payee, payeeW), truncate(tx.Narration, narrW))
			if i == m.cursor {
				sb.WriteString(cursorStyle.Render("> ") + selectedStyle.Render(row) + "\n")
			} else {
				sb.WriteString("  " + dimStyle.Render(row) + "\n")
			}
		}
	}

	sb.WriteString("\n")
	if m.loading {
		sb.WriteString(dimStyle.Render("  Refreshing...") + "\n")
	} else if m.err != nil {
		sb.WriteString(errorStyle.Render("  Error: "+m.err.Error()) + "\n")
	} else if m.status != "" {
		sb.WriteString(successStyle.Render("  "+m.status) + "\n")
	} else {
		fmt.Fprintf(&sb, "%s\n", dimStyle.Render(fmt.Sprintf("  %d transactions", len(m.items))))
	}
	sb.WriteString(dimStyle.Render("  ↑↓/jk navigate   enter detail   r refresh   q quit"))

	return sb.String()
}

func (m TransactionListModel) detailView() string {
	if len(m.items) == 0 {
		return ""
	}
	tx := m.items[m.cursor]
	var sb strings.Builder

	title := tx.Date
	if tx.Payee != "" {
		title += "  " + tx.Payee + " · " + tx.Narration
	} else if tx.Narration != "" {
		title += "  " + tx.Narration
	}
	sb.WriteString(boldStyle.Render(title) + "\n\n")

	if tx.Flag != "" {
		fmt.Fprintf(&sb, "  %-14s%s\n", "Flag:", tx.Flag)
	}
	if len(tx.Tags) > 0 {
		fmt.Fprintf(&sb, "  %-14s%s\n", "Tags:", strings.Join(tx.Tags, ", "))
	}
	if len(tx.Links) > 0 {
		fmt.Fprintf(&sb, "  %-14s%s\n", "Links:", strings.Join(tx.Links, ", "))
	}

	sb.WriteString("\n  " + boldStyle.Render("Postings:") + "\n")
	for _, p := range tx.Postings {
		fmt.Fprintf(&sb, "    %-40s %s %s\n", p.Account, p.Units.Number, p.Units.Currency)
	}

	sb.WriteString("\n")
	fmt.Fprintf(&sb, "  %-14s%s\n", "Entry hash:", dimStyle.Render(tx.EntryHash))

	if m.err != nil {
		sb.WriteString("\n" + errorStyle.Render("  Error: "+m.err.Error()) + "\n")
	}

	sb.WriteString("\n" + dimStyle.Render("  d delete   esc back   q quit"))
	return sb.String()
}

func (m TransactionListModel) confirmView() string {
	if len(m.items) == 0 {
		return ""
	}
	tx := m.items[m.cursor]
	var sb strings.Builder

	title := tx.Date
	if tx.Payee != "" {
		title += "  " + tx.Payee + " · " + tx.Narration
	} else if tx.Narration != "" {
		title += "  " + tx.Narration
	}
	sb.WriteString(boldStyle.Render(title) + "\n\n")

	sb.WriteString("\n  " + boldStyle.Render("Postings:") + "\n")
	for _, p := range tx.Postings {
		fmt.Fprintf(&sb, "    %-40s %s %s\n", p.Account, p.Units.Number, p.Units.Currency)
	}

	sb.WriteString("\n")
	sb.WriteString(warnStyle.Render("  ⚠  Delete this transaction? [y/n]") + "\n")
	sb.WriteString(dimStyle.Render("  y confirm   n/esc cancel"))
	return sb.String()
}

func (m TransactionListModel) deletingView() string {
	return dimStyle.Render("  Deleting transaction...") + "\n"
}

func truncate(s string, max int) string {
	if max <= 0 {
		return ""
	}
	if len(s) <= max {
		return s
	}
	if max <= 3 {
		return s[:max]
	}
	return s[:max-3] + "..."
}
