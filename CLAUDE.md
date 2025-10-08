# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

**mdv** is a dual-mode markdown viewer written in Go that operates as both a terminal UI (TUI) and a native GUI application. It uses Bubble Tea for the TUI, Wails for the GUI, and Glamour/Goldmark for markdown rendering.

## Build & Development Commands

All development is managed through [Task](https://taskfile.dev/). Common commands:

```bash
# Build both TUI and GUI binaries
task build:all

# Build TUI only (outputs ./mdv)
task build:tui

# Build GUI only (requires Wails CLI, outputs to cmd/mdv-gui/build/bin/mdv-gui)
task build:gui

# Run TUI in development with a file
task run:tui -- examples/demo.md

# Run GUI in development (Wails dev mode with hot reload)
task run:gui -- examples/demo.md

# Test with built binaries (simulates installed state)
task demo -- mdv examples/
task demo -- mdv -g examples/demo.md

# Run tests
task test

# Format code
task fmt

# Install to $GOPATH/bin
task install:all

# Clean build artifacts
task clean
```

**Manual Go commands** (without Task):
- Build TUI: `go build -o ./mdv ./cmd/mdv`
- Build GUI: `cd cmd/mdv-gui && wails build`
- Run TUI: `go run ./cmd/mdv [file.md]`

## Architecture

### Core Components

1. **cmd/mdv/** - TUI application entry point
   - Cobra-based CLI with smart file detection
   - Bubble Tea model/view/update loop (cmd/mdv/main.go:19-101)
   - File picker for directories with multiple markdown files (cmd/mdv/main.go:103-148)
   - File watcher for live reload using fsnotify (cmd/mdv/main.go:363-398)
   - Can launch GUI mode via `-g` flag by exec'ing mdv-gui

2. **cmd/mdv-gui/** - GUI application (Wails-based)
   - Simple Go backend (app.go) that loads markdown and converts to HTML
   - Frontend renders HTML in a native webview
   - Opens external links in system browser (cmd/mdv-gui/app.go:63-65)

3. **internal/config/** - Configuration system
   - Uses Viper for layered config (flags > directory .mdv.yaml > local .mdv.yaml > ~/.config/mdv/config.yaml > env vars)
   - Directory-specific configs: when viewing a file, mdv checks for .mdv.yaml in that file's directory
   - Environment variables: MDV_THEME, MDV_THEME_LIGHT, MDV_THEME_DARK, MDV_WRAP, MDV_WATCH, MDV_GUI, MDV_EXCLUDE, MDV_EDITOR
   - Config struct: Theme, ThemeLight, ThemeDark, Wrap, GUI, Watch, Exclude, File, Editor, GoldmarkExtensions (plugin metadata)
   - Setting `gui: true` in config makes `mdv` launch in GUI mode by default (equivalent to `-g` flag)

4. **internal/render/** - Markdown rendering
   - `ToANSI()`: Converts markdown to ANSI for terminal (uses Glamour)
   - `ToHTML()`: Converts markdown to HTML for GUI (uses Goldmark)
   - Goldmark configured with GitHub Flavored Markdown extensions and can load user-provided plugins from config
   - Theme auto-detection: checks macOS AppleInterfaceStyle or Linux COLORFGBG (internal/render/render.go:33-63)

### Key Features

- **Smart file detection**: No args → auto-detect single .md in current dir, or show picker if multiple
- **Directory scanning**: Point to a directory, mdv finds all markdown files (respects exclude patterns)
- **Exclude patterns**: Glob patterns via config or `--exclude` flag to skip files when scanning
- **Live reload**: `--watch` flag uses fsnotify to auto-update TUI on file changes
- **Dual mode**: TUI (`mdv`) can launch GUI (`mdv-gui`) by pressing 'o' or using `-g` flag
- **Theme handling**: "auto" theme detects system dark/light preference; users can customize which themes are used via `theme-light` and `theme-dark` config options

### Configuration Precedence

1. Command-line flags (highest)
2. `.mdv.yaml` in target file's directory
3. `.mdv.yaml` in current working directory
4. `~/.config/mdv/config.yaml`
5. Environment variables
6. Built-in defaults (lowest)

### TUI Keybindings

- `↑`/`↓` or `j`/`k`: Scroll
- `g`/`G` or `Home`/`End`: Jump to top/bottom
- `r`: Manual reload
- `o`: Open in GUI mode
- `q`/`Esc`/`Ctrl+C`: Quit

## Testing

The repository includes example markdown files in `examples/` for testing various features:
- `examples/demo.md` - General demo
- `examples/features.md` - Feature showcase
- `examples/config.yaml` - Example config

## Dependencies

Key external dependencies:
- **Bubble Tea**: TUI framework
- **Glamour**: Terminal markdown rendering
- **Goldmark**: Markdown parsing (with GFM extensions)
- **Wails v2**: GUI framework
- **Cobra**: CLI framework
- **Viper**: Configuration management
- **fsnotify**: File watching

## Notes

- GUI builds require Wails CLI and CGO_ENABLED=1
- The TUI can run without Wails/GUI dependencies
- Both binaries (mdv, mdv-gui) are independent executables
- mdv-gui can be invoked directly or via `mdv -g`
