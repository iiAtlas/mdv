package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestGUICommandReportsDecodedConfig(t *testing.T) {
	t.Setenv("MDV_GUI_TEST", "1")

	tmpDir := t.TempDir()
	mdFile := filepath.Join(tmpDir, "doc.md")
	if err := os.WriteFile(mdFile, []byte("# Doc"), 0o644); err != nil {
		t.Fatalf("failed to write markdown file: %v", err)
	}

	var buf bytes.Buffer
	rootCmd.SetOut(&buf)
	rootCmd.SetErr(&buf)
	rootCmd.SetArgs([]string{"--gui-theme=dark", "--gui-width=wide", mdFile})
	defer func() {
		rootCmd.SetOut(os.Stdout)
		rootCmd.SetErr(os.Stderr)
		rootCmd.SetArgs(nil)
	}()

	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("expected gui root command to succeed in test mode, got error: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "gui-theme=dark") {
		t.Fatalf("expected gui-theme to report dark, got %q", out)
	}
	if !strings.Contains(out, "gui-width=wide") {
		t.Fatalf("expected gui-width to report wide, got %q", out)
	}
	if !strings.Contains(out, "file="+mdFile) {
		t.Fatalf("expected file path to be reported, got %q", out)
	}
}

func TestGUIHelpDisplaysUsage(t *testing.T) {
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

	if !strings.Contains(buf.String(), "mdv-gui [file.md]") {
		t.Fatalf("expected help output to include usage header, got %q", buf.String())
	}
}
