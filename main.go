package main

import (
	"io/fs"
	"log"
	"net/http"
	"os"
)

const storageFolder = "storage"
const postFolder = "posts"
const aboutFolder = "about"
const aboutFileName = "about.md"

var PORT = os.Getenv("PORT")
var defaultPort = "8787"

func main() {
	if PORT == "" {
		PORT = defaultPort
	}
	
	// create storage Folder as filesystem
	filesystem := os.DirFS(storageFolder)
	
	// create post Folder as sub filesystem of storage
	postFs, err := fs.Sub(filesystem, postFolder)
	if err != nil {
		log.Fatal(err.Error())
	}

	// create about Folder as sub filesystem of storage
	aboutFs, err := fs.Sub(filesystem, aboutFolder)
	if err != nil {
		log.Fatal(err.Error())
	}

	posts, err := NewPostsFromFS(postFs)
	if err != nil {
		log.Fatal("failure when loading post folder: ", err.Error())
	}

	about, err := fs.ReadFile(aboutFs, aboutFileName)
	if err != nil {
		log.Fatal("failure when read about file: ", err.Error())
	}

	render, err := NewRender()
	if err != nil {
		log.Fatal("Failure when creating a render: ", err.Error())
	}

	http.Handle("/", indexHandler(posts, render))
	http.Handle("/post/", http.StripPrefix("/post/", postHandler(posts, render)))
	http.Handle("/about", aboutHandler(about, render))

	http.Handle("/static/",
		http.StripPrefix("/static/",
			http.FileServer(http.Dir("public"))))

	log.Println("server started at", PORT)

	if err := http.ListenAndServe(":"+PORT, nil); err != nil {
		log.Fatal("Server error: ", err.Error())
	}
}

func indexHandler(posts []Post, render *Render) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		if err := render.IndexPage(rw, posts); err != nil {
			log.Println("Failure when rendering the index: ", err.Error())
			http.Error(rw, "internal server error", http.StatusInternalServerError)
		}
	}
}

func aboutHandler(about []byte, render *Render) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		err := render.AboutPage(rw, about)
		if err != nil {
			log.Println("Failure when rendering about page: ", err.Error())
			http.Error(rw, "internal server error", http.StatusInternalServerError)
			return
		}
	}
}

func postHandler(posts []Post, render *Render) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		for _, post := range posts {
			if post.Slug == r.URL.Path {
				if err := render.PostPage(rw, post); err != nil {
					log.Println("Failure when rendering post page:", err.Error())
					http.Error(rw, "internal server error", http.StatusInternalServerError)
					return
				}
				return
			}
		}

		http.Error(rw, "Not Found", http.StatusNotFound)
	}
}
