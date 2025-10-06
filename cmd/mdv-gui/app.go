package main

import (
	"context"
	"os"

	"github.com/iiatlas/mdv/internal/render"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type App struct {
	ctx  context.Context
	html string
	path string
}

func NewApp(path string) *App { return &App{path: path} }

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	runtime.LogDebugf(ctx, "Startup called with path: '%s'", a.path)
	// If no CLI arg, ask user to pick a file
	if a.path == "" {
		runtime.LogDebug(ctx, "No path provided, opening file dialog")
		p, err := runtime.OpenFileDialog(a.ctx, runtime.OpenDialogOptions{
			Title: "Open Markdown",
			Filters: []runtime.FileFilter{{DisplayName: "Markdown", Pattern: "*.md;*.markdown"}},
		})
		if err == nil && p != "" {
			a.path = p
		}
	}
	// Load + render if we have a path
	if a.path != "" {
		runtime.LogDebugf(ctx, "Loading file: %s", a.path)
		err := a.load(a.path)
		if err != nil {
			runtime.LogErrorf(ctx, "Failed to load file: %v", err)
		} else {
			runtime.LogDebugf(ctx, "Successfully loaded file, HTML length: %d", len(a.html))
		}
	}
}

func (a *App) load(path string) error {
	b, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	htmlBytes, err := render.ToHTML(b)
	if err != nil {
		return err
	}
	a.html = string(htmlBytes)
	runtime.WindowSetTitle(a.ctx, "mdv — "+path)
	return nil
}

// GetHTML is exposed to JS
func (a *App) GetHTML() string { return a.html }

// OpenURL opens a URL in the default browser
func (a *App) OpenURL(url string) {
	runtime.BrowserOpenURL(a.ctx, url)
}
