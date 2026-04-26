package utils

import (
	"bytes"
	"strings"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
)

var blogMarkdown = goldmark.New(
	goldmark.WithExtensions(extension.GFM),
	goldmark.WithParserOptions(
		parser.WithAutoHeadingID(),
	),
	goldmark.WithRendererOptions(
		html.WithHardWraps(),
		html.WithUnsafe(),
	),
)

func MarkdownToHTML(markdown string) string {
	if strings.TrimSpace(markdown) == "" {
		return ""
	}

	var buf bytes.Buffer
	if err := blogMarkdown.Convert([]byte(markdown), &buf); err != nil {
		return markdown
	}
	return buf.String()
}
