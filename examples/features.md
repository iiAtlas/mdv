# MDV Features Guide

This document showcases the features and capabilities of **MDV** (Markdown Viewer).

## Installation & Usage

### Basic Usage

```bash
# View a markdown file
mdv document.md

# Use a specific theme
mdv -t light document.md

# Enable watch mode (auto-reload on file changes)
mdv --watch document.md

# Open in GUI mode
mdv-gui document.md
```

## Themes

MDV supports multiple rendering themes:

- **dark** (default) - Optimized for dark terminals
- **light** - Optimized for light terminals
- **auto** - Automatically detects terminal theme

### Changing Themes

```bash
# Command-line flag
mdv -t dark document.md

# Configuration file
# Create .mdv.yaml in your project or home directory
theme: "dark"
```

## Watch Mode

Watch mode automatically reloads the document when the file changes:

```bash
mdv --watch document.md
```

**Perfect for:**
- Live preview while editing
- Documentation development
- Real-time content updates

**Tip**: Edit this file in your favorite editor and run `mdv --watch features.md` to see live updates!

## Keyboard Shortcuts

| Key | Action |
|-----|--------|
| `↑` / `k` | Scroll up |
| `↓` / `j` | Scroll down |
| `g` / `Home` | Jump to top |
| `G` / `End` | Jump to bottom |
| `r` | Reload file |
| `q` / `Esc` / `Ctrl+C` | Quit |

## Configuration

MDV can be configured using a YAML file in multiple locations:

1. `.mdv.yaml` in current directory
2. `.mdv.yaml` in home directory
3. Command-line flags (override config files)

### Configuration Options

```yaml
# .mdv.yaml
theme: "dark"        # dark, light, or auto
wrap: 80            # Text wrap width
watch: false        # Enable watch mode
gui: false          # Launch in GUI mode
```

### Example Configuration

Create `.mdv.yaml` in your project:

```yaml
# Project-specific MDV config
theme: "dark"
wrap: 100
watch: true
```

## Dual Mode: TUI & GUI

MDV offers both terminal and graphical interfaces:

### TUI Mode (default)
```bash
mdv document.md
```
- Fast and lightweight
- Works over SSH
- Perfect for quick viewing

### GUI Mode
```bash
mdv-gui document.md
```
- Native window application
- Web-based rendering
- Better image support

## Advanced Features

### GitHub Flavored Markdown (GFM)

MDV fully supports GFM extensions:

#### Tables

| Feature | TUI | GUI |
|---------|-----|-----|
| Markdown rendering | ✅ | ✅ |
| Syntax highlighting | ✅ | ✅ |
| Live reload | ✅ | ✅ |
| Images | Limited | ✅ |

#### Task Lists

- [x] Basic markdown support
- [x] GFM extensions
- [x] Syntax highlighting
- [x] Multiple themes
- [ ] Plugin system (planned)

#### Strikethrough

~~This text is crossed out~~

#### Autolinks

URLs like https://github.com are automatically linked.

### Syntax Highlighting

Code blocks support syntax highlighting for many languages:

```rust
fn fibonacci(n: u32) -> u32 {
    match n {
        0 => 0,
        1 => 1,
        _ => fibonacci(n - 1) + fibonacci(n - 2),
    }
}
```

```typescript
interface User {
    id: number;
    name: string;
    email: string;
}

const users: User[] = [];
```

## Performance

MDV is designed to be fast and efficient:

- **Instant startup** - No loading delays
- **Low memory footprint** - Efficient rendering
- **Large file support** - Handles big documents smoothly

## Tips & Tricks

### Quick Preview

Create a shell alias for faster access:

```bash
# Add to ~/.bashrc or ~/.zshrc
alias md='mdv'
alias mdw='mdv --watch'
```

### Default Theme

Set your preferred theme globally:

```bash
# ~/.mdv.yaml
theme: "auto"
```

### Project Documentation

Use MDV to preview README files:

```bash
mdv README.md
```

---

## Getting Help

- Press `?` (planned feature) for help overlay
- Use `mdv --help` for command-line options
- Visit the [GitHub repository](https://github.com/iiatlas/mdv) for documentation

## What's Next?

Try these examples:

1. **Basic markdown**: `mdv examples/demo.md`
2. **This file with watch mode**: `mdv --watch examples/features.md`
3. **Different theme**: `mdv -t light examples/demo.md`
4. **GUI mode**: `mdv-gui examples/features.md`

Happy viewing! 🚀
