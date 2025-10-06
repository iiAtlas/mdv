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
	Theme string // "dark", "light", "auto", or custom
	Wrap  int    // wrap width for terminal rendering
	GUI   bool   // open GUI (Wails) instead of TUI
	Watch bool   // auto-reload on file change
	File  string // markdown file path (positional arg)
}

// NewViper sets up Viper with sensible defaults and search paths.
func NewViper() *viper.Viper {
	v := viper.New()

	// 1) Built-in defaults
	v.SetDefault("theme", "auto")
	v.SetDefault("wrap", 80)
	v.SetDefault("gui", false)
	v.SetDefault("watch", false)

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
	if err := bind("wrap", "wrap"); err != nil {
		return err
	}
	if err := bind("gui", "gui"); err != nil {
		return err
	}
	if err := bind("watch", "watch"); err != nil {
		return err
	}
	return nil
}

// Decode pulls values from Viper into a typed Config.
func Decode(v *viper.Viper, fileArg string) (Config, error) {
	cfg := Config{
		Theme: v.GetString("theme"),
		Wrap:  v.GetInt("wrap"),
		GUI:   v.GetBool("gui"),
		Watch: v.GetBool("watch"),
		File:  fileArg,
	}
	if cfg.Wrap < 0 {
		return cfg, fmt.Errorf("wrap must be >= 0")
	}
	return cfg, nil
}
