package v1

import (
	"net/http"
)

// Router ...
type Router struct {
	postController    PostController
	versionController VersionController
}

// NewRouter ...
func NewRouter(
	postController PostController,
	versionController VersionController,
) *Router {
	return &Router{
		postController,
		versionController,
	}
}

func (router *Router) ServeHTTP(w http.ResponseWriter, request *http.Request) {
	serveMux := http.NewServeMux()
	serveMux.Handle("/posts", router.postController.GetPosts())
	serveMux.Handle("/version", router.versionController.HandlerFunc())
	serveMux.ServeHTTP(w, request)
}
