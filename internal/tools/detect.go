package tools

import (
	"fmt"
	"os/exec"
	"strings"
)

// installGuide returns a tailored installation guide for the given Python pkg.
// If uv is detected on the system it is listed as the primary method.
func installGuide(pkg string) string {
	var b strings.Builder

	_, uvErr := exec.LookPath("uv")
	_, pipErr := exec.LookPath("pip3")
	if pipErr != nil {
		_, pipErr = exec.LookPath("pip")
	}

	fmt.Fprintf(&b, "Install %s by following these steps:\n\n", pkg)

	if uvErr == nil {
		// uv is present — put it first as the preferred method.
		b.WriteString("Option 1: uv (detected on your system — recommended)\n")
		fmt.Fprintf(&b, "  uv tool install %s\n\n", pkg)
		b.WriteString("Option 2: pip\n")
		if pipErr == nil {
			fmt.Fprintf(&b, "  pip install %s\n\n", pkg)
		} else {
			fmt.Fprintf(&b, "  pip install %s   # install pip first if needed\n\n", pkg)
		}
	} else if pipErr == nil {
		// pip is present, uv is not.
		b.WriteString("Option 1: pip (detected on your system)\n")
		fmt.Fprintf(&b, "  pip install %s\n\n", pkg)
		b.WriteString("Option 2: uv\n")
		b.WriteString("  # Install uv: https://docs.astral.sh/uv/getting-started/installation/\n")
		fmt.Fprintf(&b, "  uv tool install %s\n\n", pkg)
	} else {
		// Neither found — show both with setup hints.
		b.WriteString("Option 1: pip\n")
		fmt.Fprintf(&b, "  pip install %s\n", pkg)
		b.WriteString("  # If pip is not installed: https://pip.pypa.io/en/stable/installation/\n\n")
		b.WriteString("Option 2: uv (fast Python package manager)\n")
		b.WriteString("  # Install uv: https://docs.astral.sh/uv/getting-started/installation/\n")
		fmt.Fprintf(&b, "  uv tool install %s\n\n", pkg)
	}

	b.WriteString("After installation, make sure the beancount tools are on your PATH.\n")
	fmt.Fprintf(&b, "Verify with: %s --version", pkg)

	return b.String()
}

// requireTool checks that name exists in PATH.
// pkg is the Python package that provides the tool and is used in the install guide.
// If name is missing, it returns an error with a step-by-step installation guide.
func requireTool(name, pkg string) error {
	if _, err := exec.LookPath(name); err != nil {
		return fmt.Errorf("%q not found in PATH\n\n%s", name, installGuide(pkg))
	}
	return nil
}
