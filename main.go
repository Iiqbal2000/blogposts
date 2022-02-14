package main

import (
	"log"
	"net/http"
	"os"
)

const articleFolder = "storage/posts"
const aboutFile = "storage/about/about.md"

func main() {
	postFs := os.DirFS(articleFolder)

	posts, err := NewPostsFromFS(postFs)
	if err != nil {
		log.Fatal("failure when loading folder")
	}

	about, err := os.ReadFile(aboutFile)
	if err != nil {
		log.Fatal(err.Error())
	}

	renderer, err := NewRenderer()
	if err != nil {
		log.Fatal("Failure when initialization renderer", err)
	}

	http.HandleFunc("/", func(rw http.ResponseWriter, r *http.Request) {
		if err := renderer.RenderIndex(rw, posts); err != nil {
			log.Fatal("Failure when rendering the index: ", err.Error())
		}
	})

	http.Handle("/post/", http.StripPrefix("/post/", PostHandler(posts)))
	http.Handle("/about", aboutHandler(about))

	http.Handle("/static/",
		http.StripPrefix("/static/",
			http.FileServer(http.Dir("public"))))

	log.Println("server started at localhost:8787")

	http.ListenAndServe(":8787", nil)
}

func aboutHandler(about []byte) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		renderer, err := NewRenderer()
		if err != nil {
			log.Fatal("Failure when initialization renderer")
		}

		err = renderer.RenderAbout(rw, about)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func PostHandler(posts []Post) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		renderer, err := NewRenderer()
		if err != nil {
			log.Fatal("Failure when initialization renderer")
		}

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
