package render

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestResolveThemeRespectsOverrides(t *testing.T) {
	t.Setenv("COLORFGBG", "15;0") // Dark background

	theme := ResolveTheme("auto", "light-custom", "dark-custom")
	if theme != "dark-custom" {
		t.Fatalf("expected dark custom theme, got %q", theme)
	}

	t.Setenv("COLORFGBG", "0;15") // Light background
	theme = ResolveTheme("auto", "light-custom", "dark-custom")
	if theme != "light-custom" {
		t.Fatalf("expected light custom theme, got %q", theme)
	}

	theme = ResolveTheme("dark", "light-custom", "dark-custom")
	if theme != "dark" {
		t.Fatalf("expected explicit theme 'dark', got %q", theme)
	}
}

func TestLoadThemeCSSBuiltin(t *testing.T) {
	css, err := loadThemeCSS("light")
	if err != nil {
		t.Fatalf("expected to load built-in theme, got error: %v", err)
	}
	if !strings.Contains(css, ".markdown-body") {
		t.Fatalf("expected CSS to contain .markdown-body")
	}
}

func TestLoadThemeCSSCustomFileAndTildeExpansion(t *testing.T) {
	tempDir := t.TempDir()
	themePath := filepath.Join(tempDir, "custom.css")
	if err := os.WriteFile(themePath, []byte(".markdown-body{color:red;}"), 0o644); err != nil {
		t.Fatalf("failed to write custom theme: %v", err)
	}

	t.Setenv("HOME", tempDir)

	css, err := loadThemeCSS("~/custom.css")
	if err != nil {
		t.Fatalf("expected to load custom theme with tilde, got error: %v", err)
	}
	if !strings.Contains(css, "color:red") {
		t.Fatalf("expected CSS content to be returned, got %q", css)
	}
}

func TestEmbeddedThemesAreAvailable(t *testing.T) {
	for name := range themeMapping {
		name := name
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			css, err := loadThemeCSS(name)
			if err != nil {
				t.Fatalf("expected embedded theme %s to load, got error: %v", name, err)
			}
			if !strings.Contains(css, ".markdown-body") {
				t.Fatalf("expected embedded theme %s CSS to include markdown-body selector", name)
			}
		})
	}
}

func TestLoadThemeCSSMissingFile(t *testing.T) {
	if _, err := loadThemeCSS("/does/not/exist.css"); err == nil {
		t.Fatalf("expected error when custom theme missing")
	}
}

func TestExtractBackgroundColor(t *testing.T) {
	css := `.markdown-body { background-color: #0d1117; color: #fff; }`
	if got := extractBackgroundColor(css); got != "#0d1117" {
		t.Fatalf("expected background color '#0d1117', got %q", got)
	}

	css = `.markdown-body { color: #fff; }`
	if got := extractBackgroundColor(css); got != "" {
		t.Fatalf("expected empty background color, got %q", got)
	}
}

func TestGetMaxWidth(t *testing.T) {
	tests := map[string]string{
		"narrow":  "680px",
		"medium":  "900px",
		"wide":    "1200px",
		"full":    "",
		"800":     "800px",
		"unknown": "900px",
	}

	for input, expected := range tests {
		if got := getMaxWidth(input); got != expected {
			t.Fatalf("expected width %q for %q, got %q", expected, input, got)
		}
	}
}

func TestToHTMLWrapsOutputWithThemeAndWidth(t *testing.T) {
	t.Setenv("COLORFGBG", "15;0")

	src := []byte("# Title\n\nContent")
	html, err := ToHTML(src, "auto", "", "", "narrow")
	if err != nil {
		t.Fatalf("ToHTML returned error: %v", err)
	}

	htmlStr := string(html)
	if !strings.Contains(htmlStr, "<div class=\"markdown-body\">") {
		t.Fatalf("expected markdown body wrapper, got %q", htmlStr)
	}

	css, err := loadThemeCSS("dark")
	if err != nil {
		t.Fatalf("failed to load dark theme for comparison: %v", err)
	}
	bgColor := extractBackgroundColor(css)
	if bgColor != "" && !strings.Contains(htmlStr, bgColor) {
		t.Fatalf("expected background color %q in output", bgColor)
	}
	if !strings.Contains(htmlStr, "max-width: 680px") {
		t.Fatalf("expected narrow width CSS, got %q", htmlStr)
	}
}

func TestToANSIProducesRenderableOutput(t *testing.T) {
	t.Setenv("COLORFGBG", "15;0")

	src := []byte("# Heading")
	output, err := ToANSI(src, "auto", "", "", 80)
	if err != nil {
		t.Fatalf("ToANSI returned error: %v", err)
	}
	if !strings.Contains(output, "Heading") {
		t.Fatalf("expected rendered output to contain heading text, got %q", output)
	}
}
