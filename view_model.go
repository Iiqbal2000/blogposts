package main

import (
	"html/template"

	"github.com/gomarkdown/markdown"
	"github.com/microcosm-cc/bluemonday"
)

type viewModel struct {
	HTMLBody template.HTML
}

func newVM(r *Renderer, body []byte) *viewModel {
	rawHtml := markdown.ToHTML([]byte(body), r.mdParser, nil)
	sanitized := bluemonday.UGCPolicy().SanitizeBytes(rawHtml)
	return &viewModel{HTMLBody: template.HTML(sanitized)}
}
