package render

import (
	"embed"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

//go:embed themes/css/*.css
var themesFS embed.FS

// Built-in theme names that map to embedded CSS files
const (
	ThemeAuto  = "auto"
	ThemeLight = "light"
	ThemeDark  = "dark"
)

// themeMapping maps theme names to their CSS file names
var themeMapping = map[string]string{
	"auto":  "themes/css/github-markdown-auto.css",
	"light": "themes/css/github-markdown-light.css",
	"dark":  "themes/css/github-markdown-dark.css",
}

// loadThemeCSS loads CSS content for a given theme.
// If theme is a built-in name (auto/light/dark), loads from embedded FS.
// If theme is a file path, reads from filesystem.
// Returns the CSS content and any error.
func loadThemeCSS(theme string) (string, error) {
	// Expand ~ in theme path for convenience
	if strings.HasPrefix(theme, "~/") {
		if home, err := os.UserHomeDir(); err == nil {
			theme = filepath.Join(home, theme[2:])
		}
	}

	// Check if it's a built-in theme
	if cssPath, ok := themeMapping[theme]; ok {
		data, err := themesFS.ReadFile(cssPath)
		if err != nil {
			return "", fmt.Errorf("failed to load built-in theme %q: %w", theme, err)
		}
		return string(data), nil
	}

	// Try to load from filesystem (custom theme)
	data, err := os.ReadFile(theme)
	if err != nil {
		return "", fmt.Errorf("failed to load custom theme from %q: %w", theme, err)
	}
	return string(data), nil
}

// extractBackgroundColor extracts the background-color from .markdown-body CSS rule.
// Returns the color value (e.g., "#0d1117") or empty string if not found.
func extractBackgroundColor(css string) string {
	// Match .markdown-body { ... background-color: #xxx; ... }
	// This regex looks for background-color property within .markdown-body rule
	re := regexp.MustCompile(`\.markdown-body\s*\{[^}]*background-color:\s*([^;}\s]+)`)
	matches := re.FindStringSubmatch(css)
	if len(matches) >= 2 {
		return strings.TrimSpace(matches[1])
	}
	return ""
}

// getMaxWidth returns the CSS max-width value for a given width setting.
// Returns empty string for "full" (no constraint).
func getMaxWidth(width string) string {
	switch width {
	case "narrow":
		return "680px" // ~65 characters, optimal for reading
	case "medium":
		return "900px" // balanced, accommodates some tables
	case "wide":
		return "1200px" // wide tables and code blocks
	case "full":
		return "" // no constraint
	default:
		// If it looks like a number, treat as pixels (e.g., "800" -> "800px")
		// Otherwise, use default "medium"
		if regexp.MustCompile(`^\d+$`).MatchString(width) {
			return width + "px"
		}
		return "900px" // default to medium
	}
}
