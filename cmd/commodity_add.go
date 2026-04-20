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

var commodityAddCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a commodity directive to a ledger",
	Long: `Add a commodity directive to a Beancount ledger.

Example:
  beancount-cli commodity add \
    --ledger user/mybook \
    --currency AAPL \
    --date 2024-01-01`,
	RunE: runCommodityAdd,
}

var (
	commodityLedger   string
	commodityCurrency string
	commodityDate     string
)

func init() {
	commodityCmd.AddCommand(commodityAddCmd)
	commodityAddCmd.Flags().StringVar(&commodityLedger, "ledger", "", "Ledger full name (e.g. user/mybook) (required)")
	commodityAddCmd.Flags().StringVar(&commodityCurrency, "currency", "", "Commodity currency symbol (e.g. AAPL) (required)")
	commodityAddCmd.Flags().StringVar(&commodityDate, "date", "", "Commodity date in YYYY-MM-DD format (default: today)")
	_ = commodityAddCmd.MarkFlagRequired("ledger")
	_ = commodityAddCmd.MarkFlagRequired("currency")
}

func buildCommodityInput(date, currency string) (generated.LedgerCommodityInput, error) {
	if currency == "" {
		return generated.LedgerCommodityInput{}, fmt.Errorf("currency is required")
	}
	if date == "" {
		date = time.Now().Format("2006-01-02")
	}
	return generated.LedgerCommodityInput{
		Date:     date,
		Currency: strings.ToUpper(strings.TrimSpace(currency)),
	}, nil
}

func runCommodityAdd(cmd *cobra.Command, _ []string) error {
	commodity, err := buildCommodityInput(commodityDate, commodityCurrency)
	if err != nil {
		return err
	}

	client, err := utils.NewAuthedClient()
	if err != nil {
		return err
	}

	resp, err := generated.AddEntryCommodity(context.Background(), client, utils.LedgerID(commodityLedger), commodity)
	if err != nil {
		return fmt.Errorf("failed to add commodity directive: %w", err)
	}

	result := resp.AddEntryCommodity
	if !result.Success {
		msg := "unknown error"
		if result.Message != nil && *result.Message != "" {
			msg = *result.Message
		}
		return fmt.Errorf("server rejected commodity directive: %s", msg)
	}

	fmt.Fprintf(cmd.OutOrStdout(), "Commodity directive added for %s in %s on %s\n",
		commodity.Currency, commodityLedger, commodity.Date)
	return nil
}
