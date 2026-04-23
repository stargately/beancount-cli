package transaction

import (
	"github.com/spf13/cobra"

	"beancount.io/beancount-cli/cmd/transaction/list"
)

func NewCmdTransaction() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "transaction",
		Short: "Manage transactions in a ledger",
	}
	cmd.AddCommand(NewCmdAdd())
	cmd.AddCommand(list.NewCmdList())
	return cmd
}
