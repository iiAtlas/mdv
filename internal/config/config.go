package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

// Config is what the rest of your app reads.
type Config struct {
	Theme         string   // "dark", "light", "auto", or custom (for TUI/Glamour)
	ThemeLight    string   // theme to use when system is in light mode (only applies when Theme is "auto")
	ThemeDark     string   // theme to use when system is in dark mode (only applies when Theme is "auto")
	GUITheme      string   // "dark", "light", "auto", or custom CSS file (for GUI/Goldmark)
	GUIThemeLight string   // GUI theme to use when system is in light mode (only applies when GUITheme is "auto")
	GUIThemeDark  string   // GUI theme to use when system is in dark mode (only applies when GUITheme is "auto")
	GUIWidth      string   // "narrow", "medium", "wide", "full" - content width for GUI
	Wrap          int      // wrap width for terminal rendering
	GUI           bool     // open GUI (Wails) instead of TUI
	Watch         bool     // auto-reload on file change
	Exclude       []string // glob patterns for files to exclude
	File          string   // markdown file path (positional arg)
	Editor        string   // editor command to open files (defaults to $EDITOR or "vim")
}

// NewViper sets up Viper with sensible defaults and search paths.
func NewViper() *viper.Viper {
	v := viper.New()

	// 1) Built-in defaults
	v.SetDefault("theme", "auto")
	v.SetDefault("theme-light", "")
	v.SetDefault("theme-dark", "")
	v.SetDefault("gui-theme", "auto")
	v.SetDefault("gui-theme-light", "")
	v.SetDefault("gui-theme-dark", "")
	v.SetDefault("gui-width", "medium")
	v.SetDefault("wrap", 80)
	v.SetDefault("gui", false)
	v.SetDefault("watch", false)
	v.SetDefault("exclude", []string{})

	// Default editor: check $EDITOR, fallback to "vim"
	defaultEditor := os.Getenv("EDITOR")
	if defaultEditor == "" {
		defaultEditor = "vim"
	}
	v.SetDefault("editor", defaultEditor)

	// 2) Local config: ./.mdv.yaml
	v.SetConfigName(".mdv")
	v.SetConfigType("yaml")
	v.AddConfigPath(".")
	_ = v.ReadInConfig() // ignore if missing

	// 3) Global config: $XDG_CONFIG_HOME/mdv/config.yaml (or OS equivalent)
	if cfgDir, err := os.UserConfigDir(); err == nil {
		v.SetConfigName("config")
		v.SetConfigType("yaml")
		v.AddConfigPath(filepath.Join(cfgDir, "mdv"))
		_ = v.MergeInConfig() // ignore if missing
	}

	// 4) Environment variables: MDV_THEME, MDV_WRAP, MDV_GUI, MDV_WATCH
	v.SetEnvPrefix("mdv")
	v.AutomaticEnv()

	return v
}

// BindFlags wires Cobra flags into Viper (flags win over files/env).
func BindFlags(v *viper.Viper, fs *pflag.FlagSet) error {
	bind := func(key, name string) error { return v.BindPFlag(key, fs.Lookup(name)) }
	if err := bind("theme", "theme"); err != nil {
		return err
	}
	if err := bind("gui-theme", "gui-theme"); err != nil {
		return err
	}
	if err := bind("gui-width", "gui-width"); err != nil {
		return err
	}
	if err := bind("wrap", "wrap"); err != nil {
		return err
	}
	if err := bind("gui", "gui"); err != nil {
		return err
	}
	if err := bind("watch", "watch"); err != nil {
		return err
	}
	if err := bind("exclude", "exclude"); err != nil {
		return err
	}
	if err := bind("editor", "editor"); err != nil {
		return err
	}
	return nil
}

// MergeDirectoryConfig loads a .mdv.yaml from the specified directory if it exists.
// This allows directory-specific configurations when viewing files in that directory.
func MergeDirectoryConfig(v *viper.Viper, dir string) {
	// Create a new viper instance just for the directory config
	dirConfig := viper.New()
	dirConfig.SetConfigName(".mdv")
	dirConfig.SetConfigType("yaml")
	dirConfig.AddConfigPath(dir)

	// Try to read the directory-specific config
	if err := dirConfig.ReadInConfig(); err == nil {
		// Merge it into the main viper instance
		_ = v.MergeConfigMap(dirConfig.AllSettings())
	}
}

// resolveThemePath resolves a theme path relative to the config directory.
// If the path is absolute or home-relative (~), returns it as-is.
// If the path is relative, tries to resolve it relative to configDir.
// If the resolved file exists, returns the absolute path; otherwise returns the original path.
func resolveThemePath(themePath, configDir string) string {
	if themePath == "" {
		return themePath
	}

	// If absolute path or starts with ~, return as-is (will be handled by render package)
	if filepath.IsAbs(themePath) || themePath[0] == '~' {
		return themePath
	}

	// If relative path, try to resolve relative to config directory
	resolvedPath := filepath.Join(configDir, themePath)

	// Check if the resolved path exists
	if _, err := os.Stat(resolvedPath); err == nil {
		// File exists, return absolute path
		absPath, err := filepath.Abs(resolvedPath)
		if err == nil {
			return absPath
		}
	}

	// File doesn't exist or error getting absolute path, return original
	// (might be a built-in theme name)
	return themePath
}

// Decode pulls values from Viper into a typed Config.
// configDir is the directory containing the config file being used (typically the directory of the file being viewed).
func Decode(v *viper.Viper, fileArg string, configDir string) (Config, error) {
	cfg := Config{
		Theme:         resolveThemePath(v.GetString("theme"), configDir),
		ThemeLight:    resolveThemePath(v.GetString("theme-light"), configDir),
		ThemeDark:     resolveThemePath(v.GetString("theme-dark"), configDir),
		GUITheme:      resolveThemePath(v.GetString("gui-theme"), configDir),
		GUIThemeLight: resolveThemePath(v.GetString("gui-theme-light"), configDir),
		GUIThemeDark:  resolveThemePath(v.GetString("gui-theme-dark"), configDir),
		GUIWidth:      v.GetString("gui-width"),
		Wrap:          v.GetInt("wrap"),
		GUI:           v.GetBool("gui"),
		Watch:         v.GetBool("watch"),
		Exclude:       v.GetStringSlice("exclude"),
		File:          fileArg,
		Editor:        v.GetString("editor"),
	}
	if cfg.Wrap < 0 {
		return cfg, fmt.Errorf("wrap must be >= 0")
	}
	return cfg, nil
}
