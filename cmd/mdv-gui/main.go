package main

import (
	"embed"
	"fmt"
	"os"
	"path/filepath"

	"github.com/iiatlas/mdv/internal/config"
	"github.com/spf13/cobra"
	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

//go:embed all:frontend
var assets embed.FS

// Version information (set via ldflags during build)
var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

var rootCmd = &cobra.Command{
	Use:     "mdv-gui [file.md]",
	Short:   "Markdown viewer GUI",
	Long:    `Native desktop markdown viewer with live reload support.`,
	Version: fmt.Sprintf("%s (commit: %s, built: %s)", version, commit, date),
	Args:    cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		// Initialize Viper
		v := config.NewViper()

		// Bind flags to Viper
		if err := config.BindFlags(v, cmd.Flags()); err != nil {
			return fmt.Errorf("failed to bind flags: %w", err)
		}

		// Get file path from args or empty string
		var path string
		if len(args) > 0 {
			path = args[0]
			// Convert to absolute path if it's relative
			if !filepath.IsAbs(path) {
				absPath, err := filepath.Abs(path)
				if err == nil {
					path = absPath
				}
			}
			// Merge directory config
			config.MergeDirectoryConfig(v, filepath.Dir(path))
		}

		// Decode config (pass directory of file being viewed for relative theme path resolution)
		configDir := "."
		if path != "" {
			configDir = filepath.Dir(path)
		}
		cfg, err := config.Decode(v, path, configDir)
		if err != nil {
			return fmt.Errorf("invalid config: %w", err)
		}

		if os.Getenv("MDV_GUI_TEST") == "1" {
			fmt.Fprintf(cmd.OutOrStdout(), "gui-theme=%s gui-width=%s file=%s\n", cfg.GUITheme, cfg.GUIWidth, cfg.File)
			return nil
		}

		// Create application with instance of App structure
		app := NewApp(cfg)

		// Create application with options
		err = wails.Run(&options.App{
			Title:  "mdv",
			Width:  1024,
			Height: 768,
			AssetServer: &assetserver.Options{
				Assets: assets,
			},
			BackgroundColour: &options.RGBA{R: 0, G: 0, B: 0, A: 0}, // Transparent - let CSS control background
			OnStartup:        app.startup,
			OnShutdown:       app.shutdown,
			Bind: []interface{}{
				app,
			},
		})

		if err != nil {
			return fmt.Errorf("failed to run GUI: %w", err)
		}
		return nil
	},
}

func init() {
	rootCmd.Flags().StringP("theme", "t", "auto", "Theme for TUI rendering (dark, light, auto)")
	rootCmd.Flags().String("gui-theme", "auto", "Theme for GUI rendering (light, dark, auto, or path to CSS file)")
	rootCmd.Flags().String("gui-width", "medium", "Content width for GUI (narrow, medium, wide, full, or pixel value)")
	rootCmd.Flags().IntP("wrap", "w", 80, "Wrap width (not used in GUI)")
	rootCmd.Flags().Bool("watch", false, "Auto-reload on file change")
	rootCmd.Flags().StringSliceP("exclude", "e", []string{}, "Glob patterns for files to exclude (comma-separated)")
	rootCmd.Flags().BoolP("gui", "g", false, "GUI mode (always true for mdv-gui)")
	rootCmd.Flags().String("editor", "", "Editor command to open files (defaults to $EDITOR or vim)")
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
