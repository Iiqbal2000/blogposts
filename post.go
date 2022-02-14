package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log"
	"strings"
	"time"
)

const (
	titleSeparator       = "Title: "
	descriptionSeparator = "Description: "
	tagsSeparator        = "Tags: "
	dateSeparator        = "Date: "
)

type Post struct {
	Title, Description, Body, Slug, Date string
	Tags                                 []string
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
		Date:        restructureDate(readMetaLine(dateSeparator)),
		Body:        readBodyLine(scanner),
	}

	post.setSlug()

	return post, nil
}

func (p *Post) setSlug() {
	slug := new(strings.Builder)

	fmt.Fprintf(slug, "%s", strings.ToLower(strings.Replace(p.Title, " ", "-", -1)))

	p.Slug = slug.String()
}

func restructureDate(date string) string {
	var result strings.Builder
	layoutFormat := "02/01/2006"
	d, err := time.Parse(layoutFormat, date)
	if err != nil {
		log.Fatal("failure when restructuring the date", err.Error())
	}

	fmt.Fprintf(&result, "%d", d.Day())
	result.WriteString(" ")
	result.WriteString(d.Month().String())
	result.WriteString(" ")
	fmt.Fprintf(&result, "%d", d.Year())
	return result.String()
}

func readBodyLine(scanner *bufio.Scanner) string {
	scanner.Scan()

	buf := bytes.Buffer{}
	for scanner.Scan() {
		fmt.Fprintln(&buf, scanner.Text())
	}

	return strings.TrimSuffix(buf.String(), "\n")
}
