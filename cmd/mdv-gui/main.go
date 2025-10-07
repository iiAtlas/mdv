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

var rootCmd = &cobra.Command{
	Use:   "mdv-gui [file.md]",
	Short: "Markdown viewer GUI",
	Long:  `Native desktop markdown viewer with live reload support.`,
	Args:  cobra.MaximumNArgs(1),
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

		// Decode config
		cfg, err := config.Decode(v, path)
		if err != nil {
			return fmt.Errorf("invalid config: %w", err)
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
			BackgroundColour: &options.RGBA{R: 255, G: 255, B: 255, A: 255},
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
	rootCmd.Flags().StringP("theme", "t", "auto", "Theme for rendering (dark, light, auto)")
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
