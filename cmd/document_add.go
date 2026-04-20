package cmd

import (
	"context"
	"fmt"
	"time"

	"github.com/spf13/cobra"

	"beancount.io/beancount-cli/generated"
	"beancount.io/beancount-cli/internal/utils"
)

var documentAddCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a document directive to a ledger",
	Long: `Add a document directive to a Beancount ledger.

Example:
  beancount-cli document add \
    --ledger user/mybook \
    --account Assets:Bank \
    --filename receipts/2024-01-01.pdf \
    --date 2024-01-01`,
	RunE: runDocumentAdd,
}

var (
	documentLedger   string
	documentAccount  string
	documentFilename string
	documentDate     string
	documentTags     []string
	documentLinks    []string
)

func init() {
	documentCmd.AddCommand(documentAddCmd)
	documentAddCmd.Flags().StringVar(&documentLedger, "ledger", "", "Ledger full name (e.g. user/mybook) (required)")
	documentAddCmd.Flags().StringVar(&documentAccount, "account", "", "Account name (e.g. Assets:Bank) (required)")
	documentAddCmd.Flags().StringVar(&documentFilename, "filename", "", "Document filename or path (required)")
	documentAddCmd.Flags().StringVar(&documentDate, "date", "", "Document date in YYYY-MM-DD format (default: today)")
	documentAddCmd.Flags().StringArrayVar(&documentTags, "tag", nil, "Tag (repeatable)")
	documentAddCmd.Flags().StringArrayVar(&documentLinks, "link", nil, "Link (repeatable)")
	_ = documentAddCmd.MarkFlagRequired("ledger")
	_ = documentAddCmd.MarkFlagRequired("account")
	_ = documentAddCmd.MarkFlagRequired("filename")
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

func runDocumentAdd(cmd *cobra.Command, _ []string) error {
	document, err := buildDocumentInput(documentDate, documentAccount, documentFilename, documentTags, documentLinks)
	if err != nil {
		return err
	}

	client, err := utils.NewAuthedClient()
	if err != nil {
		return err
	}

	resp, err := generated.AddEntryDocument(context.Background(), client, utils.LedgerID(documentLedger), document)
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
		document.Account, documentLedger, document.Date, document.Filename)
	return nil
}
