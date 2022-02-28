package main

import (
	"log"
	"net/http"
	"os"
)

const postFolder = "storage/posts"
const aboutFile = "storage/about/about.md"

var PORT = os.Getenv("PORT")
var defaultPort = "8787"
func main() {
	if PORT == "" {
		PORT = defaultPort
	}

	postFs := os.DirFS(postFolder)

	posts, err := NewPostsFromFS(postFs)
	if err != nil {
		log.Fatal("failure when loading post folder: ", err.Error())
	}

	about, err := os.ReadFile(aboutFile)
	if err != nil {
		log.Fatal("failure when read about file: ", err.Error())
	}

	render, err := NewRender()
	if err != nil {
		log.Fatal("Failure when creating a render: ", err.Error())
	}

	http.Handle("/", homeHandler(posts, render))
	http.Handle("/post/", http.StripPrefix("/post/", postHandler(posts, render)))
	http.Handle("/about", aboutHandler(about, render))

	http.Handle("/static/",
		http.StripPrefix("/static/",
			http.FileServer(http.Dir("public"))))

	log.Println("server started at", PORT)

	if err := http.ListenAndServe(":" + PORT, nil); err != nil {
		log.Fatal("Server error: ", err.Error())
	}
}

func homeHandler(posts []Post, render *Render) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		if err := render.IndexPage(rw, posts); err != nil {
			log.Fatal("Failure when rendering the index: ", err.Error())
		}
	}
}

func aboutHandler(about []byte, render *Render) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		err := render.AboutPage(rw, about)
		if err != nil {
			http.Error(rw, "internal server error", http.StatusInternalServerError)
			log.Fatal("Failure when rendering about page: ", err.Error())
			return
		}
	}
}

func postHandler(posts []Post, render *Render) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		for _, post := range posts {
			if post.Slug == r.URL.Path {
				if err := render.PostPage(rw, post); err != nil {
					log.Fatal("Failure when rendering")
				}
				return
			}
		}

		rw.WriteHeader(http.StatusNotFound)
		rw.Write([]byte("Not Found"))
	}
}
