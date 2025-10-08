package config

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func withTempWorkingDir(t *testing.T) func() {
	t.Helper()
	tempDir := t.TempDir()
	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get working directory: %v", err)
	}
	if err := os.Chdir(tempDir); err != nil {
		t.Fatalf("failed to change working directory: %v", err)
	}
	return func() {
		if err := os.Chdir(wd); err != nil {
			t.Fatalf("failed to restore working directory: %v", err)
		}
	}
}

func TestNewViperDefaults(t *testing.T) {
	restore := withTempWorkingDir(t)
	defer restore()

	t.Setenv("EDITOR", "")
	configHome := filepath.Join(t.TempDir(), "config-home")
	t.Setenv("XDG_CONFIG_HOME", configHome)

	v := NewViper()

	if got := v.GetString("theme"); got != "auto" {
		t.Fatalf("expected default theme 'auto', got %q", got)
	}
	if got := v.GetString("gui-theme"); got != "auto" {
		t.Fatalf("expected default gui-theme 'auto', got %q", got)
	}
	if got := v.GetString("gui-width"); got != "medium" {
		t.Fatalf("expected default gui-width 'medium', got %q", got)
	}
	if got := v.GetInt("wrap"); got != 80 {
		t.Fatalf("expected default wrap 80, got %d", got)
	}
	if got := v.GetBool("gui"); got {
		t.Fatalf("expected default gui false, got true")
	}
	if got := v.GetBool("watch"); got {
		t.Fatalf("expected default watch false, got true")
	}
	if got := v.GetStringSlice("exclude"); len(got) != 0 {
		t.Fatalf("expected default exclude to be empty, got %v", got)
	}
	if got := v.GetString("editor"); got != "vim" {
		t.Fatalf("expected default editor 'vim', got %q", got)
	}
}

func TestNewViperEnvironmentOverrides(t *testing.T) {
	restore := withTempWorkingDir(t)
	defer restore()

	t.Setenv("EDITOR", "micro")
	t.Setenv("XDG_CONFIG_HOME", filepath.Join(t.TempDir(), "xdg"))
	t.Setenv("MDV_THEME", "dark")
	t.Setenv("MDV_WRAP", "120")
	t.Setenv("MDV_GUI", "true")
	t.Setenv("MDV_WATCH", "true")
	t.Setenv("MDV_EXCLUDE", "vendor,.git")
	t.Setenv("MDV_EDITOR", "nano")

	v := NewViper()

	if got := v.GetString("theme"); got != "dark" {
		t.Fatalf("expected theme from env 'dark', got %q", got)
	}
	if got := v.GetInt("wrap"); got != 120 {
		t.Fatalf("expected wrap from env 120, got %d", got)
	}
	if got := v.GetBool("gui"); !got {
		t.Fatalf("expected gui from env true, got false")
	}
	if got := v.GetBool("watch"); !got {
		t.Fatalf("expected watch from env true, got false")
	}
	if got := v.GetStringSlice("exclude"); strings.Join(got, ",") != "vendor,.git" {
		t.Fatalf("expected exclude slice from env to contain vendor and .git, got %v", got)
	}
	if got := v.GetString("editor"); got != "nano" {
		t.Fatalf("expected editor from env 'nano', got %q", got)
	}
}

func TestBindFlagsPrecedence(t *testing.T) {
	v := NewViper()

	fs := pflag.NewFlagSet("mdv", pflag.ContinueOnError)
	fs.String("theme", "auto", "")
	fs.String("gui-theme", "auto", "")
	fs.String("gui-width", "medium", "")
	fs.Int("wrap", 80, "")
	fs.Bool("gui", false, "")
	fs.Bool("watch", false, "")
	fs.StringSlice("exclude", []string{}, "")
	fs.String("editor", "vim", "")

	args := []string{
		"--theme=light",
		"--gui-theme=dark",
		"--gui-width=wide",
		"--wrap=100",
		"--gui",
		"--watch",
		"--exclude=vendor",
		"--exclude=.git",
		"--editor=micro",
	}
	if err := fs.Parse(args); err != nil {
		t.Fatalf("failed to parse flags: %v", err)
	}

	if err := BindFlags(v, fs); err != nil {
		t.Fatalf("BindFlags returned error: %v", err)
	}

	if got := v.GetString("theme"); got != "light" {
		t.Fatalf("expected theme from flags 'light', got %q", got)
	}
	if got := v.GetString("gui-theme"); got != "dark" {
		t.Fatalf("expected gui-theme from flags 'dark', got %q", got)
	}
	if got := v.GetString("gui-width"); got != "wide" {
		t.Fatalf("expected gui-width from flags 'wide', got %q", got)
	}
	if got := v.GetInt("wrap"); got != 100 {
		t.Fatalf("expected wrap from flags 100, got %d", got)
	}
	if got := v.GetBool("gui"); !got {
		t.Fatalf("expected gui from flags true, got false")
	}
	if got := v.GetBool("watch"); !got {
		t.Fatalf("expected watch from flags true, got false")
	}
	if got := v.GetStringSlice("exclude"); len(got) != 2 || got[0] != "vendor" || got[1] != ".git" {
		t.Fatalf("expected exclude slice from flags, got %v", got)
	}
	if got := v.GetString("editor"); got != "micro" {
		t.Fatalf("expected editor from flags 'micro', got %q", got)
	}
}

func TestMergeDirectoryConfig(t *testing.T) {
	restore := withTempWorkingDir(t)
	defer restore()

	t.Setenv("XDG_CONFIG_HOME", filepath.Join(t.TempDir(), "xdg"))

	v := NewViper()

	dir := t.TempDir()
	configPath := filepath.Join(dir, ".mdv.yaml")
	content := []byte("theme: dark\nwrap: 120\nexclude:\n  - vendor\n")
	if err := os.WriteFile(configPath, content, 0o644); err != nil {
		t.Fatalf("failed to write directory config: %v", err)
	}

	MergeDirectoryConfig(v, dir)

	if got := v.GetString("theme"); got != "dark" {
		t.Fatalf("expected merged theme 'dark', got %q", got)
	}
	if got := v.GetInt("wrap"); got != 120 {
		t.Fatalf("expected merged wrap 120, got %d", got)
	}
	if got := v.GetStringSlice("exclude"); len(got) != 1 || got[0] != "vendor" {
		t.Fatalf("expected merged exclude ['vendor'], got %v", got)
	}
}

func TestResolveThemePath(t *testing.T) {
	dir := t.TempDir()
	themeFile := filepath.Join(dir, "custom.json")
	if err := os.WriteFile(themeFile, []byte("{}"), 0o644); err != nil {
		t.Fatalf("failed to write theme file: %v", err)
	}

	rel := "custom.json"
	resolved := resolveThemePath(rel, dir)
	if resolved != themeFile {
		t.Fatalf("expected relative path to resolve to %q, got %q", themeFile, resolved)
	}

	nonExisting := resolveThemePath("missing.json", dir)
	if nonExisting != "missing.json" {
		t.Fatalf("expected missing relative path to stay unchanged, got %q", nonExisting)
	}

	abs := resolveThemePath(themeFile, dir)
	if abs != themeFile {
		t.Fatalf("expected absolute path to remain unchanged, got %q", abs)
	}

	tildePath := "~/custom.json"
	resolvedTilde := resolveThemePath(tildePath, dir)
	if resolvedTilde != tildePath {
		t.Fatalf("expected tilde path to remain unchanged, got %q", resolvedTilde)
	}
}

func TestDecode(t *testing.T) {
	dir := t.TempDir()
	themeLight := filepath.Join(dir, "light.json")
	if err := os.WriteFile(themeLight, []byte("{}"), 0o644); err != nil {
		t.Fatalf("failed to write light theme file: %v", err)
	}
	themeDark := filepath.Join(dir, "dark.json")
	if err := os.WriteFile(themeDark, []byte("{}"), 0o644); err != nil {
		t.Fatalf("failed to write dark theme file: %v", err)
	}
	guiTheme := filepath.Join(dir, "theme.css")
	if err := os.WriteFile(guiTheme, []byte("body{}"), 0o644); err != nil {
		t.Fatalf("failed to write gui theme file: %v", err)
	}

	v := viper.New()
	v.Set("theme", "auto")
	v.Set("theme-light", filepath.Base(themeLight))
	v.Set("theme-dark", filepath.Base(themeDark))
	v.Set("gui-theme", filepath.Base(guiTheme))
	v.Set("gui-theme-light", "light-gui.css")
	v.Set("gui-theme-dark", "dark-gui.css")
	v.Set("gui-width", "wide")
	v.Set("wrap", 100)
	v.Set("gui", true)
	v.Set("watch", true)
	v.Set("exclude", []string{"vendor"})
	v.Set("editor", "nano")

	cfg, err := Decode(v, "notes.md", dir)
	if err != nil {
		t.Fatalf("Decode returned error: %v", err)
	}

	if cfg.Theme != "auto" {
		t.Fatalf("expected Theme 'auto', got %q", cfg.Theme)
	}
	if cfg.ThemeLight != themeLight {
		t.Fatalf("expected ThemeLight resolved to %q, got %q", themeLight, cfg.ThemeLight)
	}
	if cfg.ThemeDark != themeDark {
		t.Fatalf("expected ThemeDark resolved to %q, got %q", themeDark, cfg.ThemeDark)
	}
	expectedGUITheme, err := filepath.Abs(guiTheme)
	if err != nil {
		t.Fatalf("failed to get absolute path: %v", err)
	}
	if cfg.GUITheme != expectedGUITheme {
		t.Fatalf("expected GUITheme resolved to %q, got %q", expectedGUITheme, cfg.GUITheme)
	}
	if cfg.GUIWidth != "wide" {
		t.Fatalf("expected GUIWidth 'wide', got %q", cfg.GUIWidth)
	}
	if cfg.Wrap != 100 {
		t.Fatalf("expected Wrap 100, got %d", cfg.Wrap)
	}
	if !cfg.GUI {
		t.Fatalf("expected GUI true, got false")
	}
	if !cfg.Watch {
		t.Fatalf("expected Watch true, got false")
	}
	if len(cfg.Exclude) != 1 || cfg.Exclude[0] != "vendor" {
		t.Fatalf("expected Exclude ['vendor'], got %v", cfg.Exclude)
	}
	if cfg.File != "notes.md" {
		t.Fatalf("expected File 'notes.md', got %q", cfg.File)
	}
	if cfg.Editor != "nano" {
		t.Fatalf("expected Editor 'nano', got %q", cfg.Editor)
	}
}

func TestDecodeWrapValidation(t *testing.T) {
	v := viper.New()
	v.Set("wrap", -1)

	if _, err := Decode(v, "", ""); err == nil {
		t.Fatalf("expected error for negative wrap, got nil")
	}
}

func TestMergeDirectoryConfigMissingFileIsNoOp(t *testing.T) {
	v := NewViper()
	v.Set("theme", "dark")

	MergeDirectoryConfig(v, t.TempDir())

	if got := v.GetString("theme"); got != "dark" {
		t.Fatalf("expected theme to remain unchanged when config missing, got %q", got)
	}
}

func TestMergeDirectoryConfigInvalidYAMLIsIgnored(t *testing.T) {
	restore := withTempWorkingDir(t)
	defer restore()

	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, ".mdv.yaml"), []byte("theme: [broken"), 0o644); err != nil {
		t.Fatalf("failed to write invalid config: %v", err)
	}

	v := NewViper()
	MergeDirectoryConfig(v, dir)

	if got := v.GetString("theme"); got != "auto" {
		t.Fatalf("expected invalid config to be ignored leaving default 'auto', got %q", got)
	}
}

func TestNewViperWithInvalidWrapEnvReturnsZero(t *testing.T) {
	restore := withTempWorkingDir(t)
	defer restore()

	t.Setenv("MDV_WRAP", "not-a-number")

	v := NewViper()
	if got := v.GetInt("wrap"); got != 0 {
		t.Fatalf("expected invalid wrap env to decode as 0, got %d", got)
	}
}
