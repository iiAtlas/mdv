package render

import (
	"bytes"

	"github.com/charmbracelet/glamour"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
)

var md = goldmark.New(
	goldmark.WithExtensions(
		extension.GFM,            // tables, strikethrough, task lists
		extension.Table,          // explicit table extension (redundant, ok)
		extension.Linkify,        // autolink URLs
		extension.Strikethrough,  // ~~del~~
	),
	goldmark.WithParserOptions(
		parser.WithAutoHeadingID(),
	),
	goldmark.WithRendererOptions(
		html.WithUnsafe(), // allow inline HTML in markdown
	),
)

func ToHTML(src []byte) ([]byte, error) {
	var buf bytes.Buffer
	if err := md.Convert(src, &buf); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func ToANSI(src []byte, theme string) (string, error) {
	return glamour.Render(string(src), theme)
}
