package main

import (
	"html/template"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/parser"
	"github.com/microcosm-cc/bluemonday"
)

type viewModel struct {
	HTMLBody template.HTML
}

// translate to HTML and sanitize it
func newVM(body []byte) *viewModel {
	extensions := parser.CommonExtensions | parser.AutoHeadingIDs | parser.OrderedListStart | parser.NoEmptyLineBeforeBlock | parser.BackslashLineBreak
	p := parser.NewWithExtensions(extensions)
	rawHtml := markdown.ToHTML([]byte(body), p, nil)
	sanitized := bluemonday.UGCPolicy().SanitizeBytes(rawHtml)
	return &viewModel{HTMLBody: template.HTML(sanitized)}
}
