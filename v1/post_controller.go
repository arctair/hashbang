package v1

import (
	"encoding/json"
	"net/http"
)

// PostController ...
type PostController interface {
	GetPosts() http.Handler
	CreatePost() http.Handler
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
			rw.WriteHeader(201)
			c.postRepository.Create(
				Post{
					ImageUri: "https://images.unsplash.com/photo-1603316851229-26637b4bd1b8?ixlib=rb-1.2.1&ixid=eyJhcHBfaWQiOjEyMDd9&auto=format&fit=crop&w=1400&q=80",
					Tags: []string{
						"#windy",
						"#tdd",
					},
				},
			)
		},
	)
}
