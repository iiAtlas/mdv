package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/fsnotify/fsnotify"
	"github.com/iiatlas/mdv/internal/config"
	"github.com/iiatlas/mdv/internal/render"
	"github.com/spf13/cobra"
)

type model struct {
	viewport viewport.Model
	content  string
	ready    bool
	filePath string
	cfg      config.Config
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c", "esc":
			return m, tea.Quit
		case "g", "home":
			m.viewport.GotoTop()
			return m, nil
		case "G", "end":
			m.viewport.GotoBottom()
			return m, nil
		case "r":
			// Manual reload
			return m, reloadFileCmd(m.filePath, m.cfg)
		case "o":
			// Open in GUI
			return m, openInGUICmd(m.filePath)
		}

	case openGUIMsg:
		// GUI launched, quit TUI
		return m, tea.Quit

	case reloadMsg:
		m.content = string(msg)
		m.viewport.SetContent(m.content)
		m.viewport.GotoTop()
		return m, nil

	case tea.WindowSizeMsg:
		headerHeight := lipgloss.Height(m.headerView())
		footerHeight := lipgloss.Height(m.footerView())

		if !m.ready {
			m.viewport = viewport.New(msg.Width, msg.Height-headerHeight-footerHeight)
			m.viewport.YPosition = headerHeight
			m.viewport.SetContent(m.content)
			m.ready = true
		} else {
			m.viewport.Width = msg.Width
			m.viewport.Height = msg.Height - headerHeight - footerHeight
		}
	}

	var cmd tea.Cmd
	m.viewport, cmd = m.viewport.Update(msg)
	return m, cmd
}

func (m model) View() string {
	if !m.ready {
		return "\n  Initializing..."
	}
	return fmt.Sprintf("%s\n%s\n%s", m.headerView(), m.viewport.View(), m.footerView())
}

func (m model) headerView() string {
	title := "Markdown Viewer"
	return lipgloss.NewStyle().Bold(true).Render(title)
}

func (m model) footerView() string {
	info := lipgloss.NewStyle().Faint(true).Render(
		fmt.Sprintf("%3.f%%", m.viewport.ScrollPercent()*100),
	)
	help := lipgloss.NewStyle().Faint(true).Render(" • ↑/↓: scroll • g/G: top/bottom • r: reload • o: open in GUI • q: quit")

	line := strings.Repeat("─", max(0, m.viewport.Width-lipgloss.Width(info)-lipgloss.Width(help)))
	return lipgloss.JoinHorizontal(lipgloss.Center, line, info, help)
}

// pickerModel is a simple file picker
type pickerModel struct {
	files    []string
	cursor   int
	selected string
}

func (p pickerModel) Init() tea.Cmd {
	return nil
}

func (p pickerModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c", "esc":
			return p, tea.Quit
		case "up", "k":
			if p.cursor > 0 {
				p.cursor--
			}
		case "down", "j":
			if p.cursor < len(p.files)-1 {
				p.cursor++
			}
		case "enter":
			p.selected = p.files[p.cursor]
			return p, tea.Quit
		}
	}
	return p, nil
}

func (p pickerModel) View() string {
	s := "Select a markdown file:\n\n"
	for i, file := range p.files {
		cursor := " "
		if p.cursor == i {
			cursor = ">"
			file = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("12")).Render(file)
		}
		s += fmt.Sprintf(" %s %s\n", cursor, file)
	}
	s += "\n" + lipgloss.NewStyle().Faint(true).Render("↑/↓: navigate • enter: select • q: quit")
	return s
}

type reloadMsg string

func reloadFileCmd(path string, cfg config.Config) tea.Cmd {
	return func() tea.Msg {
		data, err := os.ReadFile(path)
		if err != nil {
			return reloadMsg(fmt.Sprintf("Error reading file: %v", err))
		}
		out, err := render.ToANSI(data, cfg.Theme)
		if err != nil {
			return reloadMsg(fmt.Sprintf("Error rendering: %v", err))
		}
		return reloadMsg(out)
	}
}

type openGUIMsg struct{}

func openInGUICmd(path string) tea.Cmd {
	return func() tea.Msg {
		// Try to find and launch mdv-gui
		guiPath, err := exec.LookPath("mdv-gui")
		if err != nil {
			// mdv-gui not found, silently ignore
			return nil
		}

		// Launch mdv-gui in the background
		cmd := exec.Command(guiPath, path)
		_ = cmd.Start()

		return openGUIMsg{}
	}
}

// findMarkdownFiles returns all markdown files in the specified directory (non-recursive)
func findMarkdownFiles(dir string, excludePatterns []string) ([]string, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	var files []string
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		name := entry.Name()
		if strings.HasSuffix(name, ".md") || strings.HasSuffix(name, ".markdown") {
			// Check if file matches any exclude pattern
			excluded := false
			for _, pattern := range excludePatterns {
				matched, err := filepath.Match(pattern, name)
				if err != nil {
					// Invalid pattern, skip it
					continue
				}
				if matched {
					excluded = true
					break
				}
			}
			if !excluded {
				files = append(files, name)
			}
		}
	}
	return files, nil
}

// isDirectory checks if the given path is a directory
func isDirectory(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return info.IsDir()
}

var rootCmd = &cobra.Command{
	Use:   "mdv [file.md|directory]",
	Short: "Markdown viewer with TUI",
	Long:  `A terminal-based markdown viewer with support for themes, auto-reload, and GUI mode.`,
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		// Initialize Viper
		v := config.NewViper()

		// Bind flags to Viper
		if err := config.BindFlags(v, cmd.Flags()); err != nil {
			return fmt.Errorf("failed to bind flags: %w", err)
		}

		// Get exclude patterns early (needed for file scanning)
		excludePatterns := v.GetStringSlice("exclude")

		// Auto-detect markdown file if no args provided, or if arg is a directory
		var fileArg string
		var scanDir string
		var needsScan bool

		if len(args) == 0 {
			// No args, scan current directory
			scanDir = "."
			needsScan = true
		} else if isDirectory(args[0]) {
			// Arg is a directory, scan it
			scanDir = args[0]
			needsScan = true
		} else {
			// Arg is a file path
			fileArg = args[0]
			needsScan = false
		}

		if needsScan {
			files, err := findMarkdownFiles(scanDir, excludePatterns)
			if err != nil {
				return fmt.Errorf("failed to scan directory: %w", err)
			}
			if len(files) == 0 {
				return fmt.Errorf("no markdown files found in %s", scanDir)
			}
			if len(files) == 1 {
				// Auto-select the only file
				fileArg = filepath.Join(scanDir, files[0])
			} else {
				// Multiple files found, show picker
				picker := pickerModel{files: files}
				p := tea.NewProgram(picker)
				finalModel, err := p.Run()
				if err != nil {
					return fmt.Errorf("error running picker: %w", err)
				}
				pickerResult := finalModel.(pickerModel)
				if pickerResult.selected == "" {
					return fmt.Errorf("no file selected")
				}
				fileArg = filepath.Join(scanDir, pickerResult.selected)
			}
		}

		// Decode config
		cfg, err := config.Decode(v, fileArg)
		if err != nil {
			return fmt.Errorf("invalid config: %w", err)
		}

		// Check if GUI mode requested
		if cfg.GUI {
			// Try to exec mdv-gui if available
			guiPath, err := exec.LookPath("mdv-gui")
			if err == nil {
				// mdv-gui found, exec it
				cmd := exec.Command(guiPath, cfg.File)
				cmd.Stdin = os.Stdin
				cmd.Stdout = os.Stdout
				cmd.Stderr = os.Stderr
				return cmd.Run()
			}
			// mdv-gui not found, print helpful message
			fmt.Println("GUI mode requested, but 'mdv-gui' not found.")
			fmt.Printf("Please install mdv-gui or run: mdv-gui %s\n", cfg.File)
			return fmt.Errorf("mdv-gui not found in PATH")
		}

		// Read file
		data, err := os.ReadFile(cfg.File)
		if err != nil {
			return fmt.Errorf("error reading file: %w", err)
		}

		// Render Markdown to ANSI with configured theme
		// Note: glamour doesn't support wrap width directly via API,
		// but we can use it for future custom rendering
		out, err := render.ToANSI(data, cfg.Theme)
		if err != nil {
			return fmt.Errorf("render error: %w", err)
		}

		// Create the Bubble Tea program
		m := model{
			content:  out,
			filePath: cfg.File,
			cfg:      cfg,
		}

		p := tea.NewProgram(
			m,
			tea.WithAltScreen(),
			tea.WithMouseCellMotion(),
		)

		// If watch mode, start file watcher in background
		if cfg.Watch {
			go watchFile(cfg.File, cfg.Theme, p)
		}

		if _, err := p.Run(); err != nil {
			return fmt.Errorf("error running program: %w", err)
		}

		return nil
	},
}

func watchFile(path, theme string, p *tea.Program) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return
	}
	defer watcher.Close()

	if err := watcher.Add(path); err != nil {
		return
	}

	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				return
			}
			if event.Op&fsnotify.Write == fsnotify.Write {
				// File was modified, trigger reload
				data, err := os.ReadFile(path)
				if err != nil {
					continue
				}
				out, err := render.ToANSI(data, theme)
				if err != nil {
					continue
				}
				p.Send(reloadMsg(out))
			}
		case _, ok := <-watcher.Errors:
			if !ok {
				return
			}
		}
	}
}

func init() {
	rootCmd.Flags().StringP("theme", "t", "auto", "Theme for rendering (dark, light, auto)")
	rootCmd.Flags().IntP("wrap", "w", 80, "Wrap width for terminal rendering")
	rootCmd.Flags().BoolP("gui", "g", false, "Open in GUI mode (use mdv-gui instead)")
	rootCmd.Flags().Bool("watch", false, "Auto-reload on file change")
	rootCmd.Flags().StringSliceP("exclude", "e", []string{}, "Glob patterns for files to exclude (comma-separated)")
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
