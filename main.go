package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {
	fs := os.DirFS("posts")

	posts, err := NewPostsFromFS(fs)
	if err != nil {
		log.Fatal("failure when loading folder")
	}

	http.HandleFunc("/", func(rw http.ResponseWriter, r *http.Request) {
		postRenderer, err := NewPostRenderer()
		if err != nil {
			log.Fatal("Failure when initialization renderer", err)
		}

		if err := postRenderer.RenderIndex(rw, posts); err != nil {
			log.Fatal("Failure when rendering")
		}
	})

	http.Handle("/post/", http.StripPrefix("/post/", PostHandler(posts)))

	http.Handle("/static/", 
        http.StripPrefix("/static/", 
            http.FileServer(http.Dir("public"))))

	fmt.Println("server started at localhost:8787")

	http.ListenAndServe(":8787", nil)
}

func PostHandler(posts []Post) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		postRenderer, err := NewPostRenderer()
		if err != nil {
			log.Fatal("Failure when initialization renderer", err)
		}
		
		for _, post := range posts {
			if post.Slug == r.URL.Path {
				if err := postRenderer.Render(rw, post); err != nil {
					log.Fatal("Failure when rendering")
				}
				return
			}
		}

		rw.WriteHeader(http.StatusNotFound)
		rw.Write([]byte("Not Found"))
	}
}
