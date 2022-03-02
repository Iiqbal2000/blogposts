package main

import (
	"bytes"
	"io"
	"testing"

	approvals "github.com/approvals/go-approval-tests"
)

func TestRender(t *testing.T) {
	var (
		aPost = Post{
			Title:       "hello world",
			Slug:        "hello-world",
			Tags:        []string{"go", "tdd"},
			Body:        "## This is a post",
		}
	)

	postRenderer, err := NewRender()
	if err != nil {
		t.Fatal(err)
	}

	t.Run("it converts a single post into HTML", func(t *testing.T) {
		buf := bytes.Buffer{}

		if err := postRenderer.PostPage(&buf, aPost); err != nil {
			t.Fatal(err)
		}

		approvals.VerifyString(t, buf.String())
	})

	t.Run("it renders an index of posts", func(t *testing.T) {
		buf := bytes.Buffer{}
		posts := []Post{{Title: "Hello World", Slug: "hello-world"}, {Title: "Hello World 2", Slug: "hello-world-2"}}

		if err := postRenderer.IndexPage(&buf, posts); err != nil {
			t.Fatal(err)
		}

		approvals.VerifyString(t, buf.String())
	})
}

func BenchmarkRender(b *testing.B) {
	var (
		aPost = Post{
			Title:       "hello world",
			Slug:        "hello-world",
			Body:        "This is a post",
			Tags:        []string{"go", "tdd"},
		}
	)

	postRenderer, err := NewRender()

	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		postRenderer.PostPage(io.Discard, aPost)
	}
}
