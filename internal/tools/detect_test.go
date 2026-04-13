package tools

import (
	"strings"
	"testing"
)

func TestRequireTool_missingBinary(t *testing.T) {
	err := requireTool("__no_such_binary__", "somepkg")
	if err == nil {
		t.Fatal("expected error for missing binary, got nil")
	}

	msg := err.Error()

	// Must name the missing tool.
	if !strings.Contains(msg, "__no_such_binary__") {
		t.Errorf("error should mention the tool name, got: %s", msg)
	}

	// Must mention at least one installation method.
	if !strings.Contains(msg, "pip") && !strings.Contains(msg, "uv") {
		t.Errorf("error should include installation instructions, got: %s", msg)
	}
}

func TestRequireTool_presentBinary(t *testing.T) {
	// "go" is always in PATH when running tests.
	if err := requireTool("go", "somepkg"); err != nil {
		t.Fatalf("expected nil for present binary, got: %v", err)
	}
}

func TestInstallGuide_beancount(t *testing.T) {
	guide := installGuide("beancount")
	if !strings.Contains(guide, "pip install beancount") {
		t.Error("guide should contain 'pip install beancount'")
	}
	if !strings.Contains(guide, "uv tool install beancount") {
		t.Error("guide should contain 'uv tool install beancount'")
	}
}

func TestInstallGuide_beanquery(t *testing.T) {
	guide := installGuide("beanquery")
	if !strings.Contains(guide, "pip install beanquery") {
		t.Error("guide should contain 'pip install beanquery'")
	}
	if !strings.Contains(guide, "uv tool install beanquery") {
		t.Error("guide should contain 'uv tool install beanquery'")
	}
	// Must NOT suggest installing the wrong package.
	if strings.Contains(guide, "pip install beancount\n") {
		t.Error("guide for beanquery should not suggest 'pip install beancount'")
	}
}
