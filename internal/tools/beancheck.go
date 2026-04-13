package tools

import (
	"context"
	"os"
	"os/exec"
)

// Check runs bean-check against filename.
// Validation errors are written directly to stderr in Emacs-compatible format
// (filename:line: error message) by bean-check itself.
//
// Returns nil when the file is clean.
// Returns *exec.ExitError when bean-check finds errors (exit code reflects severity).
// Returns a descriptive error (including install guide) when bean-check is not found.
func Check(ctx context.Context, filename string) error {
	if err := requireTool("bean-check", "beancount"); err != nil {
		return err
	}
	cmd := exec.CommandContext(ctx, "bean-check", filename)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
