//go:build windows

package render

import (
	"fmt"

	"github.com/iiatlas/mdv/internal/config"
	"github.com/yuin/goldmark"
)

func loadCustomExtensions(extConfigs []config.GoldmarkExtension) ([]goldmark.Extender, error) {
	if len(extConfigs) > 0 {
		return nil, fmt.Errorf("goldmark extensions via plugins are not supported on Windows")
	}
	return nil, nil
}
