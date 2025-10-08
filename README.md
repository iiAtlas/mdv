# mdv

> A modern markdown viewer with dual TUI and GUI modes

[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](https://opensource.org/licenses/MIT)
[![Go Version](https://img.shields.io/badge/Go-1.25+-00ADD8?logo=go)](https://go.dev/)

**mdv** is a powerful, flexible markdown viewer that works both in the terminal (TUI) and as a native desktop application (GUI). It's built on top of [Glamour](https://github.com/charmbracelet/glamour) and [Goldmark](https://github.com/yuin/goldmark). With smart file detection, automatic theme matching, and live reload capabilities, `mdv` makes reading and previewing markdown effortless _and_ glamorous!

## Features

- 🎯 **Smart File Detection** - Auto-opens single files, shows picker for multiple files
- 📁 **Directory Scanning** - Point to a directory and mdv finds all markdown files
- 📂 **Multi-File Support** - Open multiple files at once in GUI mode
- 🎨 **Automatic Theme Detection** - Matches your system's dark/light mode automatically
- 🖥️ **Dual Modes** - Use as a TUI (`mdv`) or GUI (`mdv-gui`) application
- 🔄 **Live Reload** - Watch mode automatically updates when files change
- ✨ **Beautiful Rendering** - Powered by [Glamour](https://github.com/charmbracelet/glamour) & [Goldmark](https://github.com/yuin/goldmark) for styling
- 📋 **GitHub Flavored Markdown** - Full support for tables, task lists, strikethrough, and more
- ⚙️ **Flexible Configuration** - Configure via YAML files, environment variables, or command-line flags.

## Installation

### Homebrew (macOS/Linux)

Easy mode, recommended for most:
```bash
# Install the combined TUI/GUI binary from the Homebrew tap
brew install --cask iiatlas/tap/mdv
```

If you only want the TUI or GUI, they can be installed seperately:
```bash
# Only need the TUI? Install just the CLI binary
brew install --cask iiatlas/tap/mdv-tui

# Only need the GUI? Install just the GUI binary
brew install --cask iiatlas/tap/mdv-gui
```

### Download Binary

Releases are powered by [GoReleaser](https://goreleaser.com/). You can download the latest release for your platform from the [releases page](https://github.com/iiAtlas/mdv/releases).

### Build from Source

See the [Building from Source](#building-from-source) section below.

## Basic Usage

```bash
# Auto-detect and open the only markdown file in current directory
mdv

# Open a specific file
mdv README.md

# Scan a directory and pick a file
mdv examples/

# Open in GUI mode
mdv -g README.md

# Open multiple files in GUI mode (opens separate windows)
mdv -g README.md CHANGELOG.md LICENSE

# Auto-detect file and open in GUI
mdv -g

# Watch mode - auto-reload on file changes
mdv --watch document.md

# Use light theme in TUI
mdv -t light README.md

# Use dark theme in GUI with narrow width for reading
mdv -g --gui-theme dark --gui-width narrow README.md

# GUI with wide layout for technical docs
mdv-gui --gui-theme light --gui-width wide TECHNICAL.md
```

**Keyboard Shortcuts:**

| Key | Action |
|-----|--------|
| `↑` / `k` | Scroll up |
| `↓` / `j` | Scroll down |
| `g` / `Home` | Jump to top |
| `G` / `End` | Jump to bottom |
| `r` | Reload file |
| `e` | Open in editor |
| `o` | Open in GUI mode |
| `q` / `Esc` / `Ctrl+C` | Quit |

### GUI Mode

**Multi-Select Picker (GUI Mode Only):**

When scanning a directory with `-g` flag, the picker allows multi-select:
- Use `↑`/`↓` or `j`/`k` to navigate
- Press `space` to stage files (shows `[x]`)
- Press `enter` to open all staged files in separate windows
- Or just press `enter` without staging to open the current file

The GUI automatically opens external links in your default browser.

### Directory Scanning

Point mdv at any directory and it will find all markdown files:

```bash
# If directory has one .md file, opens it automatically
mdv ~/projects/my-app/

# If directory has multiple .md files, shows a picker
mdv ~/documentation/

# Works with GUI mode too (picker supports multi-select)
mdv -g ~/notes/

# Exclude specific files when scanning
mdv --exclude "README.md,draft-*.md" ~/notes/
```

**Note:** In TUI mode, the picker is single-select. In GUI mode (`-g`), you can stage multiple files with `space` and open them all.


## Configuration

mdv can be configured through multiple sources (in order of precedence):

1. Command-line flags
2. `.mdv.yaml` in target directory (when viewing files in that directory)
3. `.mdv.yaml` in current directory
4. `~/.config/mdv/config.yaml` (global config)
5. Environment variables

### Configuration File

Create a `.mdv.yaml` file in your project or home directory:

```yaml
# TUI Theme: auto, dark, light, notty, dracula, pink, tokyo-night, or path to custom JSON theme
theme: auto

# Override TUI themes for light/dark mode when theme is "auto"
# Leave empty to use default "light" and "dark" themes
# theme-light: notty
# theme-dark: tokyo-night

# GUI Theme: auto, light, dark, or path to custom CSS file
# gui-theme: auto

# Override GUI themes for light/dark mode when gui-theme is "auto"
# gui-theme-light: light
# gui-theme-dark: dark

# Content width for GUI rendering
# Options: narrow (680px), medium (900px, default), wide (1200px), full (no constraint)
# Or specify a custom pixel value: "800"
# gui-width: medium

# Text wrap width for terminal rendering
wrap: 80

# Auto-reload on file changes
watch: false

# Open in GUI mode by default
gui: false

# Editor to use when pressing 'e' (defaults to $EDITOR or vim)
# Examples: vim, nvim, nano, code, "code --wait", "subl -w"
editor: code

# Exclude files when scanning directories (glob patterns)
exclude:
  - README.md
  - "draft-*.md"
  - "*-wip.md"
```

### Environment Variables

```bash
# TUI settings
export MDV_THEME=dark
export MDV_THEME_LIGHT=notty        # Theme to use in light mode when MDV_THEME=auto
export MDV_THEME_DARK=tokyo-night   # Theme to use in dark mode when MDV_THEME=auto
export MDV_WRAP=100

# GUI settings
export MDV_GUI_THEME=dark           # GUI theme (auto, light, dark, or CSS file path)
export MDV_GUI_THEME_LIGHT=light    # GUI theme for light mode when MDV_GUI_THEME=auto
export MDV_GUI_THEME_DARK=dark      # GUI theme for dark mode when MDV_GUI_THEME=auto
export MDV_GUI_WIDTH=medium         # GUI content width (narrow, medium, wide, full, or pixel value)

# General settings
export MDV_WATCH=true
export MDV_GUI=false
export MDV_EDITOR=code
export MDV_EXCLUDE="README.md,draft-*.md"
```

### Command-line Flags

```bash
mdv --help

Usage:
  mdv [file.md|directory...]

Flags:
      --editor string       Editor command to open files (defaults to $EDITOR or vim)
  -e, --exclude strings     Glob patterns for files to exclude (comma-separated)
  -g, --gui                 Open in GUI mode (use mdv-gui instead)
      --gui-theme string    Theme for GUI rendering (light, dark, auto, or path to CSS file) (default "auto")
      --gui-width string    Content width for GUI (narrow, medium, wide, full, or pixel value) (default "medium")
  -h, --help                help for mdv
  -t, --theme string        Theme for TUI rendering (dark, light, auto) (default "auto")
      --watch               Auto-reload on file change
  -w, --wrap int            Wrap width for terminal rendering (default 80)
```

### Available Themes

#### TUI Themes (Terminal)

- **auto** - Automatically matches your system theme (default)
- **dark** - Optimized for dark terminals
- **light** - Optimized for light terminals
- **notty** - Plain text without colors
- **dracula** - Dracula color scheme
- **pink** - Pink color scheme
- **tokyo-night** - Tokyo Night color scheme

**Custom Themes for Auto Mode:**

When using `theme: auto`, you can customize which themes are used for light and dark modes:

```yaml
theme: auto
theme-light: notty      # Use 'notty' theme when system is in light mode
theme-dark: tokyo-night # Use 'tokyo-night' theme when system is in dark mode
```

This allows you to have personalized themes that automatically switch with your system appearance while still benefiting from auto-detection.

**Custom Theme JSON Files:**

You can create your own custom themes using [Glamour's JSON theme format](https://github.com/charmbracelet/glamour/tree/master/styles). Theme values can be either built-in theme names or paths to custom JSON theme files:

```yaml
# Use a built-in theme
theme: dark

# Use an absolute path to a custom theme file
theme: /Users/me/.config/mdv/themes/custom.json

# Use a relative path (resolved relative to the config file's directory)
theme: ./themes/mytheme.json

# Use home directory expansion
theme: ~/.config/mdv/themes/ocean.json

# Works with auto mode too
theme: auto
theme-light: ~/themes/light-custom.json
theme-dark: ~/themes/dark-custom.json
```

**Note on relative paths:** When using relative paths in a config file (e.g., `.mdv.yaml`), the path is resolved relative to the directory containing that config file, not relative to your current working directory. For example, if you have `examples/.mdv.yaml` with `theme: ./ocean-theme.json`, mdv will look for `examples/ocean-theme.json`.

**Creating Custom Themes:**

Custom theme files use JSON format and define colors and styles for different markdown elements. See [Glamour's style documentation](https://github.com/charmbracelet/glamour/tree/master/styles) for the complete schema and built-in theme examples.

Example custom theme structure:

```json
{
  "document": {
    "color": "252"
  },
  "heading": {
    "color": "39",
    "bold": true
  },
  "paragraph": {},
  "code_block": {
    "color": "244",
    "background_color": "236"
  }
}
```

You can start by copying an existing theme from the [Glamour styles directory](https://github.com/charmbracelet/glamour/tree/master/styles) and modifying it to your preferences.

#### GUI Themes (Desktop Application)

The GUI mode uses CSS-based themes powered by [github-markdown-css](https://github.com/sindresorhus/github-markdown-css):

- **auto** - Automatically matches your system theme (default)
- **light** - GitHub-style light theme
- **dark** - GitHub-style dark theme

**Custom Themes for Auto Mode:**

Similar to TUI themes, you can customize which themes are used in GUI mode:

```yaml
gui-theme: auto
gui-theme-light: light  # Use light theme when system is in light mode
gui-theme-dark: dark    # Use dark theme when system is in dark mode
```

**Custom CSS Theme Files:**

You can use custom CSS files for complete control over GUI appearance:

```yaml
# Use a built-in theme
gui-theme: dark

# Use an absolute path to a custom CSS file
gui-theme: /Users/me/.config/mdv/themes/custom.css

# Use a relative path (resolved relative to the config file's directory)
gui-theme: ./themes/mytheme.css

# Use home directory expansion
gui-theme: ~/.config/mdv/themes/ocean.css

# Works with auto mode too
gui-theme: auto
gui-theme-light: ~/themes/light-custom.css
gui-theme-dark: ~/themes/dark-custom.css
```

**Creating Custom CSS Themes:**

Custom CSS themes should style the `.markdown-body` class. You can start with github-markdown-css and modify it, or create your own from scratch. The CSS should include:

- `.markdown-body` base styles (background-color, color, font-family, etc.)
- Heading styles (h1-h6)
- Paragraph, list, and text formatting
- Code block and inline code styles
- Table styles
- Link styles

Example command-line usage:

```bash
# Use light theme in GUI
mdv-gui --gui-theme light README.md

# Use custom CSS theme
mdv-gui --gui-theme ~/my-theme.css README.md
```

#### GUI Width Options

Control the content width in GUI mode for optimal reading:

- **narrow** (680px) - Optimal for reading long-form content (~65 characters per line)
- **medium** (900px, default) - Balanced layout for mixed content
- **wide** (1200px) - Accommodates wide tables and code blocks
- **full** - No width constraint, uses entire window
- **Custom pixel value** - Specify exact width (e.g., `800`)

```yaml
# In config file
gui-width: narrow

# Or via command-line
mdv-gui --gui-width narrow README.md
mdv-gui --gui-width 800 README.md
```

The width setting centers the content and maintains proper background color throughout the window.

## Building from Source

### Prerequisites

- [Go 1.25 or later](https://go.dev/dl/)
- [Task](https://taskfile.dev/) (optional, but recommended)
- [Wails CLI](https://wails.io/docs/gettingstarted/installation) (for GUI builds)

### Clone and Build

```bash
# Clone the repository
git clone https://github.com/iiAtlas/mdv.git
cd mdv

# Install dependencies
go mod download

# Build TUI binary only
go build -o mdv ./cmd/mdv

# Or build both TUI and GUI using Task
task build:all
```

### Development with Task

Task provides convenient commands for development. Use `task --list` to see all available tasks.

#### Build Commands

```bash
# Build both TUI and GUI (default)
task
task build:all

# Build TUI only
task build:tui

# Build GUI only (requires Wails CLI)
task build:gui

# Install/update dependencies
task deps
```

#### Run Commands

```bash
# Run TUI in development (use -- to pass args)
task run:tui -- examples/demo.md

# Run GUI in development with hot reload (use -- to pass args)
task run:gui -- examples/demo.md

# Demo mode - test with built binaries (simulates installed state)
task demo -- mdv examples/
task demo -- mdv -g examples/demo.md

# Run TUI with built-in watch mode
task watch -- examples/demo.md
```

#### Install Commands

```bash
# Install both TUI and GUI to $GOPATH/bin
task install:all

# Install TUI only
task install:tui

# Install GUI only
task install:gui

# Uninstall both
task uninstall:all

# Uninstall TUI only
task uninstall:tui

# Uninstall GUI only
task uninstall:gui
```

#### Testing & Quality Commands

```bash
# Run tests
task test

# Run tests with coverage report (generates coverage.html)
task test:cover

# Format code
task fmt

# Run linter (requires golangci-lint)
task lint
```

#### Cleanup & Release Commands

```bash
# Clean build artifacts
task clean

# Create release via GoReleaser (requires goreleaser)
task release

# Build snapshot release without publishing
task release:snapshot

# Show all available tasks
task help
```

### Manual Build Commands

If you prefer not to use Task:

```bash
# Build TUI
go build -o ./mdv ./cmd/mdv

# Build GUI (requires Wails CLI)
cd cmd/mdv-gui
wails build
```

### Project Structure

```
mdv/
├── cmd/
│   ├── mdv/           # TUI application
│   └── mdv-gui/       # GUI application (Wails)
├── internal/
│   ├── config/        # Configuration management
│   └── render/        # Markdown rendering engine
├── examples/          # Example markdown files and configs
├── Taskfile.yaml      # Task automation
└── .goreleaser.yaml   # Release configuration
```

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request or open an issue.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Acknowledgments

An incredible shout-out is due to the [charm_](https://charm.land/) team for all they've done to make the Command Line _glamorous_.  Thank you!

- Built with [Bubble Tea](https://github.com/charmbracelet/bubbletea) for the TUI
- Styled with [Lipgloss](https://github.com/charmbracelet/lipgloss) and [Glamour](https://github.com/charmbracelet/glamour)
- GUI powered by [Wails](https://wails.io/)
- Markdown parsing via [Goldmark](https://github.com/yuin/goldmark)

---

**Made for fun by [Atlas](https://github.com/iiAtlas)**
