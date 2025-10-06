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
	// If no CLI arg, ask user to pick a file
	if a.path == "" {
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
		_ = a.load(a.path)
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
