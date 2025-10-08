package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/viper"
)

func TestDecodeGoldmarkExtensions(t *testing.T) {
	dir := t.TempDir()

	pluginA := filepath.Join(dir, "extA.so")
	if err := os.WriteFile(pluginA, []byte("test"), 0o644); err != nil {
		t.Fatalf("failed to create pluginA: %v", err)
	}

	nestedDir := filepath.Join(dir, "nested")
	if err := os.MkdirAll(nestedDir, 0o755); err != nil {
		t.Fatalf("failed to create nested dir: %v", err)
	}
	pluginB := filepath.Join(nestedDir, "extB.so")
	if err := os.WriteFile(pluginB, []byte("test"), 0o644); err != nil {
		t.Fatalf("failed to create pluginB: %v", err)
	}

	v := viper.New()
	v.Set("goldmark", map[string]any{
		"extensions": []map[string]any{
			{"path": filepath.Base(pluginA)},
			{"path": filepath.Join("nested", "extB.so"), "symbol": "Custom"},
		},
	})

	cfg, err := Decode(v, "", dir)
	if err != nil {
		t.Fatalf("Decode returned error: %v", err)
	}

	if len(cfg.GoldmarkExtensions) != 2 {
		t.Fatalf("expected 2 extensions, got %d", len(cfg.GoldmarkExtensions))
	}

	if cfg.GoldmarkExtensions[0].Path != pluginA {
		t.Fatalf("expected first extension path %q, got %q", pluginA, cfg.GoldmarkExtensions[0].Path)
	}
	if cfg.GoldmarkExtensions[0].Symbol != "Extension" {
		t.Fatalf("expected default symbol 'Extension', got %q", cfg.GoldmarkExtensions[0].Symbol)
	}

	if cfg.GoldmarkExtensions[1].Path != pluginB {
		t.Fatalf("expected second extension path %q, got %q", pluginB, cfg.GoldmarkExtensions[1].Path)
	}
	if cfg.GoldmarkExtensions[1].Symbol != "Custom" {
		t.Fatalf("expected custom symbol 'Custom', got %q", cfg.GoldmarkExtensions[1].Symbol)
	}
}
