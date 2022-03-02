package main

import (
	"html/template"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/parser"
	"github.com/microcosm-cc/bluemonday"
)

// Content represents content that is ready to serve
type Content struct {
	Post
	HTMLBody template.HTML
}

// toHTML converts markdown to html syntax
func (c *Content) toHTML(body []byte) template.HTML {
	extensions := parser.CommonExtensions | parser.AutoHeadingIDs | parser.OrderedListStart | parser.NoEmptyLineBeforeBlock | parser.BackslashLineBreak
	p := parser.NewWithExtensions(extensions)
	rawHtml := markdown.ToHTML(body, p, nil)
	return template.HTML(rawHtml)
}

// sanitize cleans bad html syntax
func (c *Content) sanitize(html template.HTML) {
	safeHtml := bluemonday.UGCPolicy().Sanitize(string(html))
	c.HTMLBody = template.HTML(safeHtml)
}
