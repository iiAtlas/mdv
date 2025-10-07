package main

import (
	"context"
	"os"

	"github.com/fsnotify/fsnotify"
	"github.com/iiatlas/mdv/internal/config"
	"github.com/iiatlas/mdv/internal/render"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type App struct {
	ctx     context.Context
	html    string
	cfg     config.Config
	watcher *fsnotify.Watcher
}

func NewApp(cfg config.Config) *App {
	return &App{
		cfg: cfg,
	}
}

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	runtime.LogDebugf(ctx, "Startup called with path: '%s'", a.cfg.File)
	// If no CLI arg, ask user to pick a file
	if a.cfg.File == "" {
		runtime.LogDebug(ctx, "No path provided, opening file dialog")
		p, err := runtime.OpenFileDialog(a.ctx, runtime.OpenDialogOptions{
			Title:   "Open Markdown",
			Filters: []runtime.FileFilter{{DisplayName: "Markdown", Pattern: "*.md;*.markdown"}},
		})
		if err == nil && p != "" {
			a.cfg.File = p
		}
	}
	// Load + render if we have a path
	if a.cfg.File != "" {
		runtime.LogDebugf(ctx, "Loading file: %s", a.cfg.File)
		err := a.load(a.cfg.File)
		if err != nil {
			runtime.LogErrorf(ctx, "Failed to load file: %v", err)
		} else {
			runtime.LogDebugf(ctx, "Successfully loaded file, HTML length: %d", len(a.html))
		}

		// Start file watcher if watch mode is enabled
		if a.cfg.Watch {
			runtime.LogDebug(ctx, "Watch mode enabled, starting file watcher")
			go a.startWatch()
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

// Reload reloads the current file and returns the new HTML
func (a *App) Reload() (string, error) {
	if a.cfg.File == "" {
		return "", nil
	}
	err := a.load(a.cfg.File)
	if err != nil {
		return "", err
	}
	return a.html, nil
}

// startWatch starts watching the file for changes
func (a *App) startWatch() {
	var err error
	a.watcher, err = fsnotify.NewWatcher()
	if err != nil {
		runtime.LogErrorf(a.ctx, "Failed to create watcher: %v", err)
		return
	}

	err = a.watcher.Add(a.cfg.File)
	if err != nil {
		runtime.LogErrorf(a.ctx, "Failed to watch file: %v", err)
		a.watcher.Close()
		return
	}

	runtime.LogDebugf(a.ctx, "Watching file: %s", a.cfg.File)

	for {
		select {
		case event, ok := <-a.watcher.Events:
			if !ok {
				return
			}
			if event.Op&fsnotify.Write == fsnotify.Write {
				runtime.LogDebugf(a.ctx, "File modified: %s", event.Name)
				// Reload the file
				err := a.load(a.cfg.File)
				if err != nil {
					runtime.LogErrorf(a.ctx, "Failed to reload file: %v", err)
					continue
				}
				// Notify frontend that content changed
				runtime.EventsEmit(a.ctx, "file-changed")
			}
		case err, ok := <-a.watcher.Errors:
			if !ok {
				return
			}
			runtime.LogErrorf(a.ctx, "Watcher error: %v", err)
		}
	}
}

// shutdown cleans up resources
func (a *App) shutdown(ctx context.Context) {
	if a.watcher != nil {
		a.watcher.Close()
	}
}
