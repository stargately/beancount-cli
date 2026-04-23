package document

import (
	"context"
	"fmt"
	"time"

	"github.com/spf13/cobra"

	"beancount.io/beancount-cli/generated"
	"beancount.io/beancount-cli/internal/utils"
)

func NewCmdAdd() *cobra.Command {
	var ledger, account, filename, date string
	var tags, links []string

	cmd := &cobra.Command{
		Use:   "add",
		Short: "Add a document directive to a ledger",
		Long: `Add a document directive to a Beancount ledger.

Example:
  beancount-cli document add \
    --ledger user/mybook \
    --account Assets:Bank \
    --filename receipts/2024-01-01.pdf \
    --date 2024-01-01`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			document, err := buildDocumentInput(date, account, filename, tags, links)
			if err != nil {
				return err
			}
			client, err := utils.NewAuthedClient()
			if err != nil {
				return err
			}
			resp, err := generated.AddEntryDocument(context.Background(), client, utils.LedgerID(ledger), document)
			if err != nil {
				return fmt.Errorf("failed to add document directive: %w", err)
			}
			result := resp.AddEntryDocument
			if !result.Success {
				msg := "unknown error"
				if result.Message != nil && *result.Message != "" {
					msg = *result.Message
				}
				return fmt.Errorf("server rejected document directive: %s", msg)
			}
			fmt.Fprintf(cmd.OutOrStdout(), "Document directive added for %s in %s on %s: %s\n",
				document.Account, ledger, document.Date, document.Filename)
			return nil
		},
	}

	cmd.Flags().StringVar(&ledger, "ledger", "", "Ledger full name (e.g. user/mybook) (required)")
	cmd.Flags().StringVar(&account, "account", "", "Account name (e.g. Assets:Bank) (required)")
	cmd.Flags().StringVar(&filename, "filename", "", "Document filename or path (required)")
	cmd.Flags().StringVar(&date, "date", "", "Document date in YYYY-MM-DD format (default: today)")
	cmd.Flags().StringArrayVar(&tags, "tag", nil, "Tag (repeatable)")
	cmd.Flags().StringArrayVar(&links, "link", nil, "Link (repeatable)")
	_ = cmd.MarkFlagRequired("ledger")
	_ = cmd.MarkFlagRequired("account")
	_ = cmd.MarkFlagRequired("filename")
	return cmd
}

func buildDocumentInput(date, account, filename string, tags, links []string) (generated.LedgerDocumentInput, error) {
	if account == "" {
		return generated.LedgerDocumentInput{}, fmt.Errorf("account is required")
	}
	if filename == "" {
		return generated.LedgerDocumentInput{}, fmt.Errorf("filename is required")
	}
	if date == "" {
		date = time.Now().Format("2006-01-02")
	}
	if tags == nil {
		tags = []string{}
	}
	if links == nil {
		links = []string{}
	}
	return generated.LedgerDocumentInput{
		Date:     date,
		Account:  account,
		Filename: filename,
		Tags:     tags,
		Links:    links,
	}, nil
}
