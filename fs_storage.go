package main

import (
	"io/fs"
)

func NewPostsFromFS(fileSystem fs.FS) ([]Post, error) {
	dirs, err := fs.ReadDir(fileSystem, ".")
	if err != nil {
		return nil, err
	}
	var posts []Post

	for _, f := range dirs {
		post, err := getPost(fileSystem, f.Name())
		if err != nil {
			return nil, err
		}

		posts = append(posts, post)
	}

	return posts, nil
}

func getPost(fileSystem fs.FS, fileName string) (Post, error) {
	postFile, err := fileSystem.Open(fileName)
	if err != nil {
		return Post{}, err
	}

	defer postFile.Close()

	return newPost(postFile)
}
