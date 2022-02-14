package main

import (
	"embed"
	"html/template"
	"io"

	"github.com/gomarkdown/markdown/parser"
)

var (
	//go:embed "views/*"
	templates embed.FS
)

type Renderer struct {
	templ    *template.Template
	mdParser *parser.Parser
}

func NewRenderer() (*Renderer, error) {
	templ, err := template.ParseFS(templates, "views/*.gohtml")
	if err != nil {
		return nil, err
	}

	extensions := parser.CommonExtensions | parser.AutoHeadingIDs | parser.OrderedListStart | parser.NoEmptyLineBeforeBlock | parser.BackslashLineBreak
	parser := parser.NewWithExtensions(extensions)

	return &Renderer{templ, parser}, nil
}

func (pr *Renderer) RenderPost(w io.Writer, p Post) error {
	vm := newVM(pr, []byte(p.Body))
	data := struct{
		Post
		HTMLBody template.HTML
	} {
		Post: p,
		HTMLBody: vm.HTMLBody,
	}
	return pr.templ.ExecuteTemplate(w, "blog.gohtml", data)
}

func (pr *Renderer) RenderIndex(w io.Writer, posts []Post) error {
	return pr.templ.ExecuteTemplate(w, "index.gohtml", posts)
}

func (pr *Renderer) RenderAbout(w io.Writer, body []byte) error {
	vm := newVM(pr, body)
	return pr.templ.ExecuteTemplate(w, "about.gohtml", vm)
}
