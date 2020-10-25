package v1

import (
	"encoding/json"
	"net/http"
)

// PostController ...
type PostController interface {
	GetPosts() http.Handler
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
