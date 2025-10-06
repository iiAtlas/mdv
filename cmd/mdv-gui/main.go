package main

import (
	"embed"
	"log"
	"os"
	"path/filepath"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

//go:embed all:frontend
var assets embed.FS

func main() {
	// Get CLI arg (optional path to markdown file)
	var path string
	if len(os.Args) > 1 {
		path = os.Args[1]
		// Convert to absolute path if it's relative
		if !filepath.IsAbs(path) {
			absPath, err := filepath.Abs(path)
			if err == nil {
				path = absPath
			}
		}
	}

	// Create application with instance of App structure
	app := NewApp(path)

	// Create application with options
	err := wails.Run(&options.App{
		Title:  "mdv",
		Width:  1024,
		Height: 768,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 255, G: 255, B: 255, A: 255},
		OnStartup:        app.startup,
		Bind: []interface{}{
			app,
		},
	})

	if err != nil {
		log.Fatal(err)
	}
}
