package cmd

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"beancount.io/beancount-cli/generated"
	"beancount.io/beancount-cli/internal/config"
	"beancount.io/beancount-cli/internal/credentials"
	"beancount.io/beancount-cli/internal/gqlclient"
)

var transactionAddCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a transaction to a ledger",
	Long: `Add a new transaction to a Beancount ledger.

Each --posting flag accepts a comma-separated value in the format:
  account,amount,currency

Example:
  beancount-cli transaction add \
    --ledger user/mybook \
    --date 2024-01-15 \
    --payee "Starbucks" \
    --narration "Coffee" \
    --posting "Expenses:Food:Coffee,5.00,USD" \
    --posting "Assets:Checking,-5.00,USD"`,
	RunE: runTransactionAdd,
}

var (
	txLedger    string
	txDate      string
	txFlag      string
	txPayee     string
	txNarration string
	txPostings  []string
	txTags      []string
	txLinks     []string
)

func init() {
	transactionCmd.AddCommand(transactionAddCmd)
	transactionAddCmd.Flags().StringVar(&txLedger, "ledger", "", "Ledger full name (e.g. user/mybook) (required)")
	transactionAddCmd.Flags().StringVar(&txDate, "date", "", "Transaction date in YYYY-MM-DD format (default: today)")
	transactionAddCmd.Flags().StringVar(&txFlag, "flag", "*", "Transaction flag: * (cleared) or ! (pending)")
	transactionAddCmd.Flags().StringVar(&txPayee, "payee", "", "Payee (optional)")
	transactionAddCmd.Flags().StringVar(&txNarration, "narration", "", "Narration / description (optional)")
	transactionAddCmd.Flags().StringArrayVar(&txPostings, "posting", nil, "Posting in account,amount,currency format (repeatable, at least one required)")
	transactionAddCmd.Flags().StringArrayVar(&txTags, "tag", nil, "Tag (repeatable, optional)")
	transactionAddCmd.Flags().StringArrayVar(&txLinks, "link", nil, "Link (repeatable, optional)")
	_ = transactionAddCmd.MarkFlagRequired("ledger")
}

func parsePosting(s string) (generated.LedgerPostingInput, error) {
	parts := strings.SplitN(s, ",", 3)
	if len(parts) != 3 {
		return generated.LedgerPostingInput{}, fmt.Errorf("invalid posting %q: expected account,amount,currency", s)
	}
	return generated.LedgerPostingInput{
		Account: parts[0],
		Units: generated.LedgerAmountInput{
			Number:   parts[1],
			Currency: parts[2],
		},
	}, nil
}

func runTransactionAdd(cmd *cobra.Command, _ []string) error {
	if len(txPostings) == 0 {
		return fmt.Errorf("at least one --posting is required")
	}

	date := txDate
	if date == "" {
		date = time.Now().Format("2006-01-02")
	}

	postings := make([]generated.LedgerPostingInput, 0, len(txPostings))
	for _, p := range txPostings {
		posting, err := parsePosting(p)
		if err != nil {
			return err
		}
		postings = append(postings, posting)
	}

	creds, err := credentials.Load()
	if err != nil {
		return fmt.Errorf("failed to read credentials: %w", err)
	}
	if creds == nil || creds.Token == "" {
		return fmt.Errorf("not logged in. Run 'beancount-cli login' to authenticate")
	}
	if creds.IsExpired() {
		return fmt.Errorf("your session has expired. Run 'beancount-cli login' to re-authenticate")
	}

	cfg := config.Load()
	client := gqlclient.NewAuthed(cfg.APIURL, creds.Token)

	ledgersResp, err := generated.ListUserOwnedLedgers(context.Background(), client)
	if err != nil {
		return fmt.Errorf("failed to list ledgers: %w", err)
	}

	var ledgerID string
	for _, l := range ledgersResp.ListUserOwnedLedgers {
		if l.FullName == txLedger {
			ledgerID = l.Id
			break
		}
	}
	if ledgerID == "" {
		return fmt.Errorf("ledger not found: %s", txLedger)
	}

	tags := txTags
	if tags == nil {
		tags = []string{}
	}
	links := txLinks
	if links == nil {
		links = []string{}
	}

	tx := generated.LedgerTransactionInput{
		Date:     date,
		Flag:     txFlag,
		Postings: postings,
		Tags:     tags,
		Links:    links,
	}
	if txPayee != "" {
		tx.Payee = &txPayee
	}
	if txNarration != "" {
		tx.Narration = &txNarration
	}

	resp, err := generated.AddEntryTransaction(context.Background(), client, ledgerID, tx)
	if err != nil {
		return fmt.Errorf("failed to add transaction: %w", err)
	}

	result := resp.AddEntryTransaction
	if !result.Success {
		msg := "unknown error"
		if result.Message != nil && *result.Message != "" {
			msg = *result.Message
		}
		return fmt.Errorf("server rejected transaction: %s", msg)
	}

	fmt.Fprintf(cmd.OutOrStdout(), "Transaction added to %s on %s\n", txLedger, date)
	return nil
}
