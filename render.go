package main

import (
	"embed"
	"html/template"
	"io"
)

var (
	//go:embed "views/*"
	templates embed.FS
)

type Renderer struct {
	templ *template.Template
}

func NewRenderer() (*Renderer, error) {
	templ, err := template.ParseFS(templates, "views/*.gohtml")
	if err != nil {
		return nil, err
	}

	return &Renderer{templ}, nil
}

func (pr *Renderer) RenderPost(w io.Writer, p Post) error {
	vm := newVM([]byte(p.Body))

	// representing data that is served in post page
	data := struct {
		Post
		HTMLBody template.HTML
	}{
		Post:     p,
		HTMLBody: vm.HTMLBody,
	}
	return pr.templ.ExecuteTemplate(w, "blog.gohtml", data)
}

func (pr *Renderer) RenderIndex(w io.Writer, posts []Post) error {
	return pr.templ.ExecuteTemplate(w, "index.gohtml", posts)
}

func (pr *Renderer) RenderAbout(w io.Writer, body []byte) error {
	vm := newVM(body)
	return pr.templ.ExecuteTemplate(w, "about.gohtml", vm)
}
