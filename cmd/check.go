package cmd

import (
	"context"
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"

	"beancount.io/beancount-cli/internal/tools"
)

var checkCmd = &cobra.Command{
	Use:   "check <file>",
	Short: "Check a Beancount ledger file for errors",
	Long: `Check parses and validates a Beancount ledger file using bean-check.

Detected errors are printed to stderr in Emacs-compatible format:
  filename:line: Error: description

Exits with a non-zero status code when errors are found.`,
	Args: cobra.ExactArgs(1),
	RunE: runCheck,
}

func init() {
	rootCmd.AddCommand(checkCmd)
}

func runCheck(cmd *cobra.Command, args []string) error {
	err := tools.Check(context.Background(), args[0])
	if err == nil {
		fmt.Fprintf(cmd.OutOrStdout(), "No errors found in %s\n", args[0])
		return nil
	}

	// bean-check already wrote its diagnostics to stderr.
	// Exit with the same code it used rather than printing a redundant error.
	var exitErr *exec.ExitError
	if e, ok := err.(*exec.ExitError); ok {
		exitErr = e
		os.Exit(exitErr.ExitCode())
	}

	// Tool not found, file not readable, etc. — surface the message.
	return err
}
