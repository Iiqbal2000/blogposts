package main

import (
	"embed"
	"html/template"
	"io"
)

// embed files that are in the views folder to program
var (
	//go:embed "views/*"
	templates embed.FS
)

// Render represents a renderer to render a page
type Render struct {
	templ *template.Template
}

func NewRender() (*Render, error) {
	templ, err := template.ParseFS(templates, "views/*.gohtml")
	if err != nil {
		return nil, err
	}

	return &Render{templ}, nil
}

func (r *Render) PostPage(w io.Writer, p Post) error {
	content := &Content{
		Post: p,
	}

	html := content.toHTML([]byte(p.Body))
	content.sanitize(html)

	return r.templ.ExecuteTemplate(w, "blog.gohtml", content)
}

func (r *Render) IndexPage(w io.Writer, posts []Post) error {
	return r.templ.ExecuteTemplate(w, "index.gohtml", posts)
}

func (r *Render) AboutPage(w io.Writer, body []byte) error {
	content := &Content{}
	html := content.toHTML(body)
	content.sanitize(html)

	return r.templ.ExecuteTemplate(w, "about.gohtml", content)
}
