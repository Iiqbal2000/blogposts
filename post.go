package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"strings"
)

const (
	titleSeparator       = "Title: "
	descriptionSeparator = "Description: "
	tagsSeparator        = "Tags: "
)

type Post struct {
	Title, Description, Body, Slug string
	Tags                           []string
}

func newPost(postFile io.Reader) (Post, error) {
	scanner := bufio.NewScanner(postFile)

	readMetaLine := func(tagName string) string {
		scanner.Scan()
		return strings.TrimPrefix(scanner.Text(), tagName)
	}

	post := Post{
		Title:       readMetaLine(titleSeparator),
		Description: readMetaLine(descriptionSeparator),
		Tags:        strings.Split(readMetaLine(tagsSeparator), ", "),
		Body:        readBodyLine(scanner),
	}

	post.setSlug()

	return post, nil
}

func (p *Post) setSlug() {
	slug := new(strings.Builder)
	// randTime := time.Now().UnixNano()

	fmt.Fprintf(slug, "%s", strings.ToLower(strings.Replace(p.Title, " ", "-", -1)))

	p.Slug = slug.String()
}

func readBodyLine(scanner *bufio.Scanner) string {
	scanner.Scan()

	buf := bytes.Buffer{}
	for scanner.Scan() {
		fmt.Fprintln(&buf, scanner.Text())
	}

	return strings.TrimSuffix(buf.String(), "\n")
}
