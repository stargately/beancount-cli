package tools

import (
	"context"
	"os"
	"os/exec"
)

// QueryOptions controls how bean-query is invoked.
type QueryOptions struct {
	// Format is the output format: "text" (default) or "csv".
	Format string
}

// Query runs bean-query against filename.
//
// When bqlQuery is non-empty, bean-query runs in batch mode and writes results
// to stdout. When bqlQuery is empty, bean-query enters its interactive shell
// (stdin is forwarded so the user can type queries at the prompt).
//
// Returns *exec.ExitError on non-zero exit, or a descriptive error (including
// install guide) when bean-query is not found.
func Query(ctx context.Context, filename, bqlQuery string, opts QueryOptions) error {
	if err := requireTool("bean-query", "beanquery"); err != nil {
		return err
	}

	args := []string{}
	if opts.Format == "csv" {
		args = append(args, "-f", "csv")
	}
	args = append(args, filename)
	if bqlQuery != "" {
		args = append(args, bqlQuery)
	}

	cmd := exec.CommandContext(ctx, "bean-query", args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
