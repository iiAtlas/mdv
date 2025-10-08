package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestCLIFlagPrecedenceOverridesEnvironment(t *testing.T) {
	t.Setenv("MDV_CLI_TEST", "1")
	t.Setenv("MDV_THEME", "light")
	t.Setenv("MDV_WRAP", "70")
	t.Setenv("MDV_GUI", "true")

	tmpDir := t.TempDir()
	mdFile := filepath.Join(tmpDir, "sample.md")
	if err := os.WriteFile(mdFile, []byte("# Sample\n\nContent"), 0o644); err != nil {
		t.Fatalf("failed to write temporary markdown file: %v", err)
	}

	var buf bytes.Buffer
	rootCmd.SetOut(&buf)
	rootCmd.SetErr(&buf)
	rootCmd.SetArgs([]string{"--theme=dark", "--wrap=120", "--gui=false", mdFile})
	defer func() {
		rootCmd.SetOut(os.Stdout)
		rootCmd.SetErr(os.Stderr)
		rootCmd.SetArgs(nil)
	}()

	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("rootCmd.Execute returned error: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "theme=dark") {
		t.Fatalf("expected theme from flag to be reported as dark, got output %q", out)
	}
	if !strings.Contains(out, "wrap=120") {
		t.Fatalf("expected wrap from flag to be reported as 120, got output %q", out)
	}
	if !strings.Contains(out, "gui=false") {
		t.Fatalf("expected gui flag to override env to false, got output %q", out)
	}
	if !strings.Contains(out, "file="+mdFile) {
		t.Fatalf("expected output to mention file %s, got %q", mdFile, out)
	}
}

func TestCLIHelpDisplaysUsage(t *testing.T) {
	var buf bytes.Buffer
	rootCmd.SetOut(&buf)
	rootCmd.SetErr(&buf)
	rootCmd.SetArgs([]string{"--help"})
	defer func() {
		rootCmd.SetOut(os.Stdout)
		rootCmd.SetErr(os.Stderr)
		rootCmd.SetArgs(nil)
	}()

	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("expected help command to succeed, got error: %v", err)
	}

	if !strings.Contains(buf.String(), "mdv [file.md|directory...]") {
		t.Fatalf("expected help output to include usage header, got %q", buf.String())
	}
}
