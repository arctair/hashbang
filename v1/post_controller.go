package v1

import (
	"encoding/json"
	"net/http"
)

// PostController ...
type PostController interface {
	GetPosts() http.Handler
}

type postController struct{}

func (c *postController) GetPosts() http.Handler {
	return http.HandlerFunc(
		func(rw http.ResponseWriter, r *http.Request) {
			bytes, err := json.Marshal(
				[]Post{
					{
						ImageUri: "https://images.unsplash.com/photo-1603316851229-26637b4bd1b8?ixlib=rb-1.2.1&ixid=eyJhcHBfaWQiOjEyMDd9&auto=format&fit=crop&w=1400&q=80",
						Tags: []string{
							"#windy",
							"#tdd",
						},
					},
				},
			)
			if err != nil {
				panic(err)
			}
			rw.Write(bytes)
		},
	)
}

// NewPostController ...
func NewPostController() PostController {
	return &postController{}
}
