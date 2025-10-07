# mdv

> A modern markdown viewer with dual TUI and GUI modes

[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](https://opensource.org/licenses/MIT)
[![Go Version](https://img.shields.io/badge/Go-1.25+-00ADD8?logo=go)](https://go.dev/)

**mdv** is a powerful, flexible markdown viewer that works both in the terminal (TUI) and as a native desktop application (GUI). With smart file detection, automatic theme matching, and live reload capabilities, mdv makes reading and previewing markdown effortless.

## Features

- 🎯 **Smart File Detection** - Auto-opens single files, shows picker for multiple files
- 📁 **Directory Scanning** - Point to a directory and mdv finds all markdown files
- 📂 **Multi-File Support** - Open multiple files at once in GUI mode
- 🎨 **Automatic Theme Detection** - Matches your system's dark/light mode automatically
- 🖥️ **Dual Modes** - Use as a TUI (`mdv`) or GUI (`mdv-gui`) application
- 🔄 **Live Reload** - Watch mode automatically updates when files change
- ✨ **Beautiful Rendering** - Powered by [Glamour](https://github.com/charmbracelet/glamour) with syntax highlighting
- 📋 **GitHub Flavored Markdown** - Full support for tables, task lists, strikethrough, and more
- ⚙️ **Flexible Configuration** - Configure via YAML files, environment variables, or command-line flags

## Installation

### Homebrew (macOS/Linux)

```bash
# Coming soon
brew install iiatlas/tap/mdv
```

### npm/npx

```bash
# Coming soon
npm install -g @iiatlas/mdv
# or use directly
npx @iiatlas/mdv README.md
```

### Download Binary

Download the latest release for your platform from the [releases page](https://github.com/iiAtlas/mdv/releases).

### Build from Source

See the [Building from Source](#building-from-source) section below.

## Quick Start

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

# Use light theme
mdv -t light README.md
```

## Usage

### TUI Mode

The terminal interface provides a fast, keyboard-driven experience:

```bash
# Open a file
mdv document.md

# Open with watch mode for live editing
mdv --watch document.md

# Use a specific theme
mdv -t dark document.md
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

Launch the native desktop application with the `-g` flag or use `mdv-gui` directly:

```bash
# Launch GUI with a file
mdv -g README.md
mdv-gui README.md

# Open multiple files (each in a separate window)
mdv -g README.md CHANGELOG.md docs/guide.md

# Launch GUI and pick from current directory
mdv -g

# Launch GUI and pick from a specific directory
mdv -g examples/
```

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

**GUI Keyboard Shortcuts:**

| Key | Action |
|-----|--------|
| `e` | Open in editor |

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
# Theme: auto, dark, light, notty, dracula, pink, tokyo-night
theme: auto

# Override themes for light/dark mode when theme is "auto"
# Leave empty to use default "light" and "dark" themes
# theme-light: notty
# theme-dark: tokyo-night

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
export MDV_THEME=dark
export MDV_THEME_LIGHT=notty        # Theme to use in light mode when MDV_THEME=auto
export MDV_THEME_DARK=tokyo-night   # Theme to use in dark mode when MDV_THEME=auto
export MDV_WRAP=100
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
      --editor string     Editor command to open files (defaults to $EDITOR or vim)
  -e, --exclude strings   Glob patterns for files to exclude (comma-separated)
  -g, --gui               Open in GUI mode (use mdv-gui instead)
  -h, --help              help for mdv
  -t, --theme string      Theme for rendering (dark, light, auto) (default "auto")
      --watch             Auto-reload on file change
  -w, --wrap int          Wrap width for terminal rendering (default 80)
```

### Available Themes

- **auto** - Automatically matches your system theme (default)
- **dark** - Optimized for dark terminals
- **light** - Optimized for light terminals
- **notty** - Plain text without colors
- **dracula** - Dracula color scheme
- **pink** - Pink color scheme
- **tokyo-night** - Tokyo Night color scheme

#### Custom Themes for Auto Mode

When using `theme: auto`, you can customize which themes are used for light and dark modes:

```yaml
theme: auto
theme-light: notty      # Use 'notty' theme when system is in light mode
theme-dark: tokyo-night # Use 'tokyo-night' theme when system is in dark mode
```

This allows you to have personalized themes that automatically switch with your system appearance while still benefiting from auto-detection.

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

- Built with [Bubble Tea](https://github.com/charmbracelet/bubbletea) for the TUI
- Styled with [Lipgloss](https://github.com/charmbracelet/lipgloss) and [Glamour](https://github.com/charmbracelet/glamour)
- GUI powered by [Wails](https://wails.io/)
- Markdown parsing via [Goldmark](https://github.com/yuin/goldmark)

---

**Made with ❤️ by [Atlas](https://github.com/iiAtlas)**
