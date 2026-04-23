package cmd

import (
	"os"

	"github.com/spf13/cobra"

	"beancount.io/beancount-cli/internal/updater"

	cmdaccount "beancount.io/beancount-cli/cmd/account"
	cmdauth "beancount.io/beancount-cli/cmd/auth"
	cmdbalance "beancount.io/beancount-cli/cmd/balance"
	cmdbudget "beancount.io/beancount-cli/cmd/budget"
	cmdcheck "beancount.io/beancount-cli/cmd/check"
	cmdcollaborator "beancount.io/beancount-cli/cmd/collaborator"
	cmdcommodity "beancount.io/beancount-cli/cmd/commodity"
	cmddocument "beancount.io/beancount-cli/cmd/document"
	cmdevent "beancount.io/beancount-cli/cmd/event"
	cmdledger "beancount.io/beancount-cli/cmd/ledger"
	cmdnote "beancount.io/beancount-cli/cmd/note"
	cmdprice "beancount.io/beancount-cli/cmd/price"
	cmdquery "beancount.io/beancount-cli/cmd/query"
	cmdtransaction "beancount.io/beancount-cli/cmd/transaction"
	cmdupgrade "beancount.io/beancount-cli/cmd/upgrade"
)

var updateCheckCh <-chan updater.UpdateResult

var rootCmd = &cobra.Command{
	Use:   "beancount-cli",
	Short: "Official CLI for Beancount",
	// Dev-only environment variables (not shown in help output):
	//   BEANCOUNT_API_URL               Override the GraphQL endpoint (default: https://beancount.io/api-gateway/)
	//   BEANCOUNT_DASHBOARD_URL         Override the dashboard URL (default: https://beancount.io)
	//   BEANCOUNT_NO_UPDATE_NOTIFIER    Set to any value to disable update notifications
	Long: `beancount-cli is the official command-line interface for Beancount.

Use it to authenticate and interact with your beancount.io account.`,
	// Don't print usage on every RunE error — the error message is sufficient.
	SilenceUsage: true,
	PersistentPreRun: func(cmd *cobra.Command, _ []string) {
		if cmd.Name() != "upgrade" {
			updateCheckCh = updater.StartUpdateCheck(cmd.Root().Version)
		}
	},
	PersistentPostRun: func(cmd *cobra.Command, _ []string) {
		if cmd.Name() == "upgrade" || updateCheckCh == nil {
			return
		}
		select {
		case result := <-updateCheckCh:
			updater.PrintUpdateTip(os.Stderr, result)
		default:
		}
	},
}

// SetVersion sets the version string shown by --version.
func SetVersion(v string) {
	rootCmd.Version = v
}

// Execute runs the root command and exits on error.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(cmdaccount.NewCmdAccount())
	rootCmd.AddCommand(cmdcollaborator.NewCmdCollaborator())
	rootCmd.AddCommand(cmdauth.NewCmdLogin())
	rootCmd.AddCommand(cmdauth.NewCmdLogout())
	rootCmd.AddCommand(cmdauth.NewCmdWhoami())
	rootCmd.AddCommand(cmdbalance.NewCmdBalance())
	rootCmd.AddCommand(cmdbudget.NewCmdBudget())
	rootCmd.AddCommand(cmdcheck.NewCmdCheck())
	rootCmd.AddCommand(cmdcommodity.NewCmdCommodity())
	rootCmd.AddCommand(cmddocument.NewCmdDocument())
	rootCmd.AddCommand(cmdevent.NewCmdEvent())
	rootCmd.AddCommand(cmdledger.NewCmdLedger())
	rootCmd.AddCommand(cmdnote.NewCmdNote())
	rootCmd.AddCommand(cmdprice.NewCmdPrice())
	rootCmd.AddCommand(cmdquery.NewCmdBeanQuery())
	rootCmd.AddCommand(cmdquery.NewCmdQuery())
	rootCmd.AddCommand(cmdtransaction.NewCmdTransaction())
	rootCmd.AddCommand(cmdupgrade.NewCmdUpgrade(&rootCmd.Version))
}
