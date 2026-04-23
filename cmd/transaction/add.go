package transaction

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"beancount.io/beancount-cli/generated"
	"beancount.io/beancount-cli/internal/utils"
)

func NewCmdAdd() *cobra.Command {
	var ledger, date, flag, payee, narration string
	var postings, tags, links []string

	cmd := &cobra.Command{
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
		RunE: func(cmd *cobra.Command, _ []string) error {
			tx, err := buildTransactionInput(flag, date, payee, narration, postings, tags, links)
			if err != nil {
				return err
			}
			client, err := utils.NewAuthedClient()
			if err != nil {
				return err
			}
			resp, err := generated.AddEntryTransaction(context.Background(), client, utils.LedgerID(ledger), tx)
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
			fmt.Fprintf(cmd.OutOrStdout(), "Transaction added to %s on %s\n", ledger, tx.Date)
			return nil
		},
	}

	cmd.Flags().StringVar(&ledger, "ledger", "", "Ledger full name (e.g. user/mybook) (required)")
	cmd.Flags().StringVar(&date, "date", "", "Transaction date in YYYY-MM-DD format (default: today)")
	cmd.Flags().StringVar(&flag, "flag", "*", "Transaction flag: * (cleared) or ! (pending)")
	cmd.Flags().StringVar(&payee, "payee", "", "Payee (optional)")
	cmd.Flags().StringVar(&narration, "narration", "", "Narration / description (optional)")
	cmd.Flags().StringArrayVar(&postings, "posting", nil, "Posting in account,amount,currency format (repeatable, at least one required)")
	cmd.Flags().StringArrayVar(&tags, "tag", nil, "Tag (repeatable, optional)")
	cmd.Flags().StringArrayVar(&links, "link", nil, "Link (repeatable, optional)")
	_ = cmd.MarkFlagRequired("ledger")
	return cmd
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

	postingInputs := make([]generated.LedgerPostingInput, 0, len(postingStrs))
	for _, p := range postingStrs {
		posting, err := parsePosting(p)
		if err != nil {
			return generated.LedgerTransactionInput{}, err
		}
		postingInputs = append(postingInputs, posting)
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
		Postings: postingInputs,
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
