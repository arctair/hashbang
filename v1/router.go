package v1

import (
	"net/http"
)

// Router ...
type Router struct {
	namedTagListController NamedTagListController
	versionController      VersionController
}

// NewRouter ...
func NewRouter(
	namedTagListController NamedTagListController,
	versionController VersionController,
) *Router {
	return &Router{
		namedTagListController,
		versionController,
	}
}

func (router *Router) ServeHTTP(w http.ResponseWriter, request *http.Request) {
	serveMux := http.NewServeMux()
	switch request.Method {
	case http.MethodGet:
		serveMux.Handle("/namedTagLists", router.namedTagListController.GetNamedTagLists())
		serveMux.Handle("/version", router.versionController.HandlerFunc())
	case http.MethodPost:
		serveMux.Handle("/namedTagLists", router.namedTagListController.CreateNamedTagList())
	case http.MethodDelete:
		serveMux.Handle("/namedTagLists", router.namedTagListController.DeleteNamedTagLists())
	}
	serveMux.ServeHTTP(w, request)
}
