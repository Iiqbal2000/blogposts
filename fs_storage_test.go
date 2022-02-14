package main

import (
	"reflect"
	"testing"
	"testing/fstest"
)

func TestNewBlogPosts(t *testing.T) {
	const (
		firstBody = `Title: Post 1
Description: Description 1
Tags: tdd, go
---
Hello
World`
		secondBody = `Title: Post 2
Description: Description 2
Tags: rust
---
BLM`
	)

	// represent storage
	fs := fstest.MapFS{
		"hello world.md":  {Data: []byte(firstBody)},
		"hello-world2.md": {Data: []byte(secondBody)},
	}

	posts, err := NewPostsFromFS(fs)
	if err != nil {
		t.Fatal(err)
	}

	assertPost(t, posts[0], Post{
		Title:       "Post 1",
		Slug:        "post-1",
		Description: "Description 1",
		Tags:        []string{"tdd", "go"},
		Body: `Hello
World`,
	})

	assertPost(t, posts[1], Post{
		Title:       "Post 2",
		Slug:        "post-2",
		Description: "Description 2",
		Tags:        []string{"rust"},
		Body:        `BLM`,
	})
}

func assertPost(t *testing.T, got, want Post) {
	t.Helper()
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %+v, want %+v", got, want)
	}
}
