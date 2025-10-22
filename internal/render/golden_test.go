package render

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
)

func TestToHTMLMatchesGolden(t *testing.T) {
	src := []byte("# Title\n\nContent")
	themePath := filepath.Join("testdata", "html_theme.css")

	out, err := ToHTML(src, themePath, "", "", "800", "")
	if err != nil {
		t.Fatalf("ToHTML returned error: %v", err)
	}

	goldenPath := filepath.Join("testdata", "tohtml.golden")
	if os.Getenv("MDV_UPDATE_GOLDEN") == "1" {
		if err := os.WriteFile(goldenPath, out, 0o644); err != nil {
			t.Fatalf("failed to update HTML golden file: %v", err)
		}
		return
	}
	want, err := os.ReadFile(goldenPath)
	if err != nil {
		t.Fatalf("failed to read HTML golden file: %v", err)
	}

	if !bytes.Equal(out, want) {
		t.Fatalf("ToHTML output did not match golden file\n--- got ---\n%s\n--- want ---\n%s", out, want)
	}
}

func TestToANSIMatchesGolden(t *testing.T) {
	src := []byte("# Title\n\nContent")

	out, err := ToANSI(src, "dark", "", "", 80)
	if err != nil {
		t.Fatalf("ToANSI returned error: %v", err)
	}

	goldenPath := filepath.Join("testdata", "toansi.golden")
	if os.Getenv("MDV_UPDATE_GOLDEN") == "1" {
		if err := os.WriteFile(goldenPath, []byte(out), 0o644); err != nil {
			t.Fatalf("failed to update ANSI golden file: %v", err)
		}
		return
	}
	want, err := os.ReadFile(goldenPath)
	if err != nil {
		t.Fatalf("failed to read ANSI golden file: %v", err)
	}

	if out != string(want) {
		t.Fatalf("ToANSI output did not match golden file\n--- got ---\n%q\n--- want ---\n%q", out, string(want))
	}
}
