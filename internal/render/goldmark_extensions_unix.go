//go:build !windows

package render

import (
	"fmt"
	"plugin"

	"github.com/iiatlas/mdv/internal/config"
	"github.com/yuin/goldmark"
)

func loadCustomExtensions(extConfigs []config.GoldmarkExtension) ([]goldmark.Extender, error) {
	if len(extConfigs) == 0 {
		return nil, nil
	}

	extensions := make([]goldmark.Extender, 0, len(extConfigs))
	for idx, extCfg := range extConfigs {
		if extCfg.Path == "" {
			return nil, fmt.Errorf("goldmark extension %d has empty path", idx)
		}

		mod, err := plugin.Open(extCfg.Path)
		if err != nil {
			return nil, fmt.Errorf("load goldmark extension %q: %w", extCfg.Path, err)
		}

		symbolName := extCfg.Symbol
		if symbolName == "" {
			symbolName = "Extension"
		}

		sym, err := mod.Lookup(symbolName)
		if err != nil {
			return nil, fmt.Errorf("lookup symbol %q in %q: %w", symbolName, extCfg.Path, err)
		}

		switch ext := sym.(type) {
		case goldmark.Extender:
			extensions = append(extensions, ext)
		case func() goldmark.Extender:
			extensions = append(extensions, ext())
		default:
			return nil, fmt.Errorf("symbol %q in %q does not implement goldmark.Extender", symbolName, extCfg.Path)
		}
	}

	return extensions, nil
}
