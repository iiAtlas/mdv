package render

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strconv"
	"strings"

	"github.com/charmbracelet/glamour"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
)

var md = goldmark.New(
	goldmark.WithExtensions(
		extension.GFM,           // tables, strikethrough, task lists
		extension.Table,         // explicit table extension (redundant, ok)
		extension.Linkify,       // autolink URLs
		extension.Strikethrough, // ~~del~~
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
	// Respect explicit terminal background hint when available.
	if colorfgbg := os.Getenv("COLORFGBG"); colorfgbg != "" {
		parts := strings.Split(colorfgbg, ";")
		if len(parts) > 0 {
			bgPart := strings.TrimSpace(parts[len(parts)-1])
			if bg, err := strconv.Atoi(bgPart); err == nil {
				// Treat low values as dark backgrounds, higher as light.
				if bg >= 0 && bg <= 8 {
					return "dark"
				}
				return "light"
			}
		}
	}

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
				if bg, err := strconv.Atoi(strings.TrimSpace(parts[1])); err == nil {
					if bg >= 0 && bg <= 8 {
						return "dark"
					}
					return "light"
				}
			}
		}
		// Default to dark for Linux
		return "dark"
	case "windows":
		// Windows: query registry for system theme preference
		// HKEY_CURRENT_USER\Software\Microsoft\Windows\CurrentVersion\Themes\Personalize
		// AppsUseLightTheme: 0x0 = dark, 0x1 = light
		cmd := exec.Command("reg", "query",
			"HKCU\\Software\\Microsoft\\Windows\\CurrentVersion\\Themes\\Personalize",
			"/v", "AppsUseLightTheme")
		output, err := cmd.Output()
		if err == nil {
			// Parse registry output
			// Output format: "    AppsUseLightTheme    REG_DWORD    0x0" or "0x1"
			if strings.Contains(string(output), "0x0") {
				return "dark"
			}
			return "light"
		}
		// Fallback to dark if registry query fails
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

// processLocalImages processes HTML and converts relative image paths to data URIs.
// basePath is the directory of the markdown file, used to resolve relative paths.
func processLocalImages(htmlContent string, basePath string) string {
	// Regex to match img tags with src attributes
	imgRegex := regexp.MustCompile(`<img\s+([^>]*?)src=["']([^"']+)["']([^>]*?)>`)

	return imgRegex.ReplaceAllStringFunc(htmlContent, func(match string) string {
		// Extract the src attribute value
		srcRegex := regexp.MustCompile(`src=["']([^"']+)["']`)
		srcMatches := srcRegex.FindStringSubmatch(match)
		if len(srcMatches) < 2 {
			return match
		}

		src := srcMatches[1]

		// Skip if it's already a data URI or absolute URL
		if strings.HasPrefix(src, "data:") ||
		   strings.HasPrefix(src, "http://") ||
		   strings.HasPrefix(src, "https://") ||
		   strings.HasPrefix(src, "//") {
			return match
		}

		// Resolve relative path to absolute path
		var absPath string
		if filepath.IsAbs(src) {
			absPath = src
		} else {
			absPath = filepath.Join(basePath, src)
		}

		// Read the image file
		imgData, err := os.ReadFile(absPath)
		if err != nil {
			// If we can't read the file, return the original tag
			return match
		}

		// Determine MIME type based on file extension
		ext := strings.ToLower(filepath.Ext(absPath))
		mimeType := "image/png" // default
		switch ext {
		case ".jpg", ".jpeg":
			mimeType = "image/jpeg"
		case ".png":
			mimeType = "image/png"
		case ".gif":
			mimeType = "image/gif"
		case ".svg":
			mimeType = "image/svg+xml"
		case ".webp":
			mimeType = "image/webp"
		case ".bmp":
			mimeType = "image/bmp"
		case ".ico":
			mimeType = "image/x-icon"
		}

		// Encode to base64
		encoded := base64.StdEncoding.EncodeToString(imgData)
		dataURI := fmt.Sprintf("data:%s;base64,%s", mimeType, encoded)

		// Replace src attribute with data URI
		newMatch := srcRegex.ReplaceAllString(match, fmt.Sprintf(`src="%s"`, dataURI))
		return newMatch
	})
}

func ToHTML(src []byte, theme, themeLight, themeDark, width, mdFilePath string) ([]byte, error) {
	// Convert markdown to HTML
	var buf bytes.Buffer
	if err := md.Convert(src, &buf); err != nil {
		return nil, err
	}

	// Resolve theme (handles "auto" detection)
	resolvedTheme := ResolveTheme(theme, themeLight, themeDark)

	// Load CSS for the resolved theme
	css, err := loadThemeCSS(resolvedTheme)
	if err != nil {
		return nil, err
	}

	// Extract background color from the CSS and apply to body/html
	bgColor := extractBackgroundColor(css)
	var bodyCSS string
	if bgColor != "" {
		bodyCSS = fmt.Sprintf("html, body { background-color: %s; margin: 0; padding: 0; }\n", bgColor)
	}

	// Apply width constraint if specified
	maxWidth := getMaxWidth(width)
	var widthCSS string
	if maxWidth != "" {
		widthCSS = fmt.Sprintf(".markdown-body { max-width: %s; margin-left: auto !important; margin-right: auto !important; }\n", maxWidth)
	}

	// Process local images if we have a file path
	htmlContent := buf.String()
	if mdFilePath != "" {
		basePath := filepath.Dir(mdFilePath)
		htmlContent = processLocalImages(htmlContent, basePath)
	}

	// Wrap the HTML with the CSS theme and markdown-body container
	var output bytes.Buffer
	output.WriteString("<style>\n")
	output.WriteString(bodyCSS)
	output.WriteString(css)
	output.WriteString("\n")
	output.WriteString(widthCSS) // Width CSS comes last to override any theme margins
	output.WriteString("</style>\n")
	output.WriteString(`<div class="markdown-body">` + "\n")
	output.WriteString(htmlContent)
	output.WriteString("\n</div>")

	return output.Bytes(), nil
}

func ToANSI(src []byte, theme, themeLight, themeDark string, wrap int) (string, error) {
	resolvedTheme := ResolveTheme(theme, themeLight, themeDark)

	// Expand ~ in theme path for convenience
	if strings.HasPrefix(resolvedTheme, "~/") {
		if home, err := os.UserHomeDir(); err == nil {
			resolvedTheme = filepath.Join(home, resolvedTheme[2:])
		}
	}

	// Use NewTermRenderer with WithStylePath to support both built-in themes
	// and custom JSON theme files
	r, err := glamour.NewTermRenderer(
		glamour.WithStylePath(resolvedTheme),
		glamour.WithWordWrap(wrap),
	)
	if err != nil {
		return "", err
	}

	return r.Render(string(src))
}
