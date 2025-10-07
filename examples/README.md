# Configuration Examples

mdv supports flexible configuration through multiple sources with the following priority (highest to lowest):

1. **Command-line flags** (highest priority)
2. **Environment variables** (`MDV_*`)
3. **Local config** (`./.mdv.yaml` in current directory)
4. **Global config** (`~/.config/mdv/config.yaml`)
5. **Built-in defaults** (lowest priority)

## Configuration Options

| Option | Type | Default | Description |
|--------|------|---------|-------------|
| `theme` | string | `"dark"` | Rendering theme: built-in name (dark, light, auto, etc.) or path to custom JSON theme file |
| `wrap` | int | `80` | Text wrap width (0 to disable) |
| `gui` | bool | `false` | Advisory flag for GUI mode |
| `watch` | bool | `false` | Auto-reload on file changes |

## Usage Examples

### Global Configuration

Create `~/.config/mdv/config.yaml`:

```bash
mkdir -p ~/.config/mdv
cp config.yaml ~/.config/mdv/
```

Edit the file to set your preferred defaults.

### Local Configuration

For project-specific settings, create `.mdv.yaml` in your project directory:

```bash
cp .mdv.yaml /path/to/your/project/
```

This will override global settings when you run `mdv` from that directory.

### Environment Variables

Override any setting with environment variables:

```bash
export MDV_THEME=light
export MDV_WRAP=100
export MDV_WATCH=true
mdv README.md
```

### Command-line Flags

Flags have the highest priority:

```bash
mdv --theme light --wrap 100 --watch README.md
mdv -t dracula -w 120 README.md
```

## Available Themes

### Built-in Themes

- `dark` - Dark background (default)
- `light` - Light background
- `auto` - Detect terminal background
- `notty` - Plain text output
- `dracula` - Dracula color scheme
- `pink` - Pink color scheme
- `tokyo-night` - Tokyo Night color scheme

### Custom Theme Files

You can create custom themes using JSON files. The theme option accepts both built-in theme names and paths to custom JSON theme files:

```yaml
# Use a custom theme file
theme: ./ocean-theme.json

# Or use an absolute path
theme: ~/.config/mdv/themes/custom.json

# Works with auto mode too
theme: auto
theme-light: ./light-custom.json
theme-dark: ./dark-custom.json
```

Try the included example:

```bash
# Test the ocean theme example
mdv --theme ./ocean-theme.json demo.md

# Or set it in your config
echo "theme: ./ocean-theme.json" > .mdv.yaml
mdv demo.md
```

See `ocean-theme.json` in this directory for a complete example. You can copy and modify it to create your own themes. For more details on the theme format, see [Glamour's style documentation](https://github.com/charmbracelet/glamour/tree/master/styles).

## Watch Mode

Enable watch mode to auto-reload the file when it changes:

```bash
mdv --watch README.md
```

Or set in config:

```yaml
watch: true
```

Press `r` to manually reload at any time.

## GUI Mode

The `gui` flag is advisory only. To actually use GUI mode, run:

```bash
mdv-gui README.md
```
