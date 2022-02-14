package main

import (
	"log"
	"net/http"
	"os"
)

const articleFolder = "storage/posts"
const aboutFile = "storage/about/about.md"

var PORT = os.Getenv("PORT")

func main() {
	postFs := os.DirFS(articleFolder)

	posts, err := NewPostsFromFS(postFs)
	if err != nil {
		log.Fatal("failure when loading post folder: ", err.Error())
	}

	about, err := os.ReadFile(aboutFile)
	if err != nil {
		log.Fatal("failure when read about file: ", err.Error())
	}

	renderer, err := NewRenderer()
	if err != nil {
		log.Fatal("Failure when creating renderer: ", err.Error())
	}

	http.Handle("/", homeHandler(posts, renderer))
	http.Handle("/post/", http.StripPrefix("/post/", postHandler(posts, renderer)))
	http.Handle("/about", aboutHandler(about, renderer))

	http.Handle("/static/",
		http.StripPrefix("/static/",
			http.FileServer(http.Dir("public"))))

	log.Println("server started at ", PORT)

	http.ListenAndServe(PORT, nil)
}

func homeHandler(posts []Post, renderer *Renderer) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		if err := renderer.RenderIndex(rw, posts); err != nil {
			log.Fatal("Failure when rendering the index: ", err.Error())
		}
	}
}

func aboutHandler(about []byte, renderer *Renderer) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		err := renderer.RenderAbout(rw, about)
		if err != nil {
			http.Error(rw, "internal server error", http.StatusInternalServerError)
			log.Fatal("Failure when rendering about page: ", err.Error())
			return
		}
	}
}

func postHandler(posts []Post, renderer *Renderer) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		for _, post := range posts {
			if post.Slug == r.URL.Path {
				if err := renderer.RenderPost(rw, post); err != nil {
					log.Fatal("Failure when rendering")
				}
				return
			}
		}

		rw.WriteHeader(http.StatusNotFound)
		rw.Write([]byte("Not Found"))
	}
}
