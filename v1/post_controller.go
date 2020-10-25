package v1

import (
	"encoding/json"
	"net/http"
)

// PostController ...
type PostController interface {
	GetPosts() http.Handler
	CreatePost() http.Handler
	DeletePost() http.Handler
}

type postController struct {
	postRepository PostRepository
}

// NewPostController ...
func NewPostController(postRepository PostRepository) PostController {
	return &postController{
		postRepository,
	}
}

func (c *postController) GetPosts() http.Handler {
	return http.HandlerFunc(
		func(rw http.ResponseWriter, r *http.Request) {
			bytes, err := json.Marshal(
				c.postRepository.FindAll(),
			)
			if err != nil {
				panic(err)
			}
			rw.Write(bytes)
		},
	)
}

func (c *postController) CreatePost() http.Handler {
	return http.HandlerFunc(
		func(rw http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			var post Post
			if err := json.NewDecoder(r.Body).Decode(&post); err != nil {
				panic(err)
			}
			c.postRepository.Create(post)
			rw.WriteHeader(201)
		},
	)
}

func (c *postController) DeletePost() http.Handler {
	return http.HandlerFunc(
		func(rw http.ResponseWriter, r *http.Request) {
			c.postRepository.DeleteAll()
			rw.WriteHeader(204)
		},
	)
}
