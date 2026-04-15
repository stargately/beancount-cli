package cmd

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"beancount.io/beancount-cli/generated"
	"beancount.io/beancount-cli/internal/utils"
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

// buildTransactionInput constructs a LedgerTransactionInput from raw CLI values.
// Pure function — no I/O, no network calls.
func buildTransactionInput(
	flag, date, payee, narration string,
	postingStrs, tags, links []string,
) (generated.LedgerTransactionInput, error) {
	if len(postingStrs) == 0 {
		return generated.LedgerTransactionInput{}, fmt.Errorf("at least one --posting is required")
	}

	if date == "" {
		date = time.Now().Format("2006-01-02")
	}

	postings := make([]generated.LedgerPostingInput, 0, len(postingStrs))
	for _, p := range postingStrs {
		posting, err := parsePosting(p)
		if err != nil {
			return generated.LedgerTransactionInput{}, err
		}
		postings = append(postings, posting)
	}

	if tags == nil {
		tags = []string{}
	}
	if links == nil {
		links = []string{}
	}

	tx := generated.LedgerTransactionInput{
		Date:     date,
		Flag:     flag,
		Postings: postings,
		Tags:     tags,
		Links:    links,
	}
	if payee != "" {
		tx.Payee = &payee
	}
	if narration != "" {
		tx.Narration = &narration
	}

	return tx, nil
}

func runTransactionAdd(cmd *cobra.Command, _ []string) error {
	tx, err := buildTransactionInput(txFlag, txDate, txPayee, txNarration, txPostings, txTags, txLinks)
	if err != nil {
		return err
	}

	client, err := utils.NewAuthedClient()
	if err != nil {
		return err
	}

	resp, err := generated.AddEntryTransaction(context.Background(), client, utils.LedgerID(txLedger), tx)
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

	fmt.Fprintf(cmd.OutOrStdout(), "Transaction added to %s on %s\n", txLedger, tx.Date)
	return nil
}
