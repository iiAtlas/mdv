package main

import (
	"fmt"
	"os"
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
		}

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
	help := lipgloss.NewStyle().Faint(true).Render(" • ↑/↓: scroll • g/G: top/bottom • r: reload • q: quit")

	line := strings.Repeat("─", max(0, m.viewport.Width-lipgloss.Width(info)-lipgloss.Width(help)))
	return lipgloss.JoinHorizontal(lipgloss.Center, line, info, help)
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

var rootCmd = &cobra.Command{
	Use:   "mdv [file.md]",
	Short: "Markdown viewer with TUI",
	Long:  `A terminal-based markdown viewer with support for themes, auto-reload, and GUI mode.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		// Initialize Viper
		v := config.NewViper()

		// Bind flags to Viper
		if err := config.BindFlags(v, cmd.Flags()); err != nil {
			return fmt.Errorf("failed to bind flags: %w", err)
		}

		// Decode config
		cfg, err := config.Decode(v, args[0])
		if err != nil {
			return fmt.Errorf("invalid config: %w", err)
		}

		// Check if GUI mode requested
		if cfg.GUI {
			fmt.Println("GUI mode requested. Please use 'mdv-gui' command instead.")
			fmt.Printf("Example: mdv-gui %s\n", cfg.File)
			return nil
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
	rootCmd.Flags().StringP("theme", "t", "dark", "Theme for rendering (dark, light, auto)")
	rootCmd.Flags().IntP("wrap", "w", 80, "Wrap width for terminal rendering")
	rootCmd.Flags().BoolP("gui", "g", false, "Open in GUI mode (use mdv-gui instead)")
	rootCmd.Flags().Bool("watch", false, "Auto-reload on file change")
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
