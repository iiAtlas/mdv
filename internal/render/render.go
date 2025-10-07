package render

import (
	"bytes"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/charmbracelet/glamour"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
)

var md = goldmark.New(
	goldmark.WithExtensions(
		extension.GFM,            // tables, strikethrough, task lists
		extension.Table,          // explicit table extension (redundant, ok)
		extension.Linkify,        // autolink URLs
		extension.Strikethrough,  // ~~del~~
	),
	goldmark.WithParserOptions(
		parser.WithAutoHeadingID(),
	),
	goldmark.WithRendererOptions(
		html.WithUnsafe(), // allow inline HTML in markdown
	),
)

// detectSystemTheme detects the system's dark/light mode preference
func detectSystemTheme() string {
	switch runtime.GOOS {
	case "darwin":
		// macOS: check AppleInterfaceStyle
		cmd := exec.Command("defaults", "read", "-g", "AppleInterfaceStyle")
		output, err := cmd.Output()
		if err == nil && strings.TrimSpace(string(output)) == "Dark" {
			return "dark"
		}
		return "light"
	case "linux":
		// Linux: check COLORFGBG environment variable
		// Format is typically "15;0" (light fg on dark bg) or "0;15" (dark fg on light bg)
		if colorfgbg := os.Getenv("COLORFGBG"); colorfgbg != "" {
			parts := strings.Split(colorfgbg, ";")
			if len(parts) == 2 {
				// If background is dark (0-8), use dark theme
				bg := parts[1]
				if bg >= "0" && bg <= "8" {
					return "dark"
				}
				return "light"
			}
		}
		// Default to dark for Linux
		return "dark"
	default:
		// Default to dark for other platforms
		return "dark"
	}
}

// ResolveTheme resolves "auto" to the system theme, or returns the theme as-is.
// When theme is "auto", it detects the system's dark/light preference and uses
// themeLight or themeDark if they are set, otherwise falls back to "light" or "dark".
func ResolveTheme(theme, themeLight, themeDark string) string {
	if theme != "auto" {
		return theme
	}

	// Detect system theme
	detected := detectSystemTheme()

	// Use custom theme for light/dark mode if configured
	if detected == "dark" && themeDark != "" {
		return themeDark
	}
	if detected == "light" && themeLight != "" {
		return themeLight
	}

	// Fall back to detected theme
	return detected
}

func ToHTML(src []byte) ([]byte, error) {
	var buf bytes.Buffer
	if err := md.Convert(src, &buf); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func ToANSI(src []byte, theme, themeLight, themeDark string) (string, error) {
	resolvedTheme := ResolveTheme(theme, themeLight, themeDark)
	return glamour.Render(string(src), resolvedTheme)
}
