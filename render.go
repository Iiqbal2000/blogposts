package main

import (
	"embed"
	"html/template"
	"io"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/parser"
)

var (
	//go:embed "views/*"
	postTemplate embed.FS
)

type PostRenderer struct {
	templ    *template.Template
	mdParser *parser.Parser
}

func NewPostRenderer() (*PostRenderer, error) {
	templ, err := template.ParseFS(postTemplate, "views/*.gohtml")
	if err != nil {
		return nil, err
	}

	extensions := parser.CommonExtensions | parser.AutoHeadingIDs
	parser := parser.NewWithExtensions(extensions)

	return &PostRenderer{templ, parser}, nil
}

func (pr *PostRenderer) Render(w io.Writer, p Post) error {
	vm := newPostVM(p, pr)
	return pr.templ.ExecuteTemplate(w, "blog.gohtml", vm)
}

func (pr *PostRenderer) RenderIndex(w io.Writer, posts []Post) error {
	return pr.templ.ExecuteTemplate(w, "index.gohtml", posts)
}

type postViewModel struct {
	Post
	HTMLBody template.HTML
}

func newPostVM(p Post, r *PostRenderer) postViewModel {
	vm := postViewModel{Post: p}

	vm.HTMLBody = template.HTML(markdown.ToHTML([]byte(p.Body), r.mdParser, nil))
	return vm
}
