package v1

import (
	"encoding/json"
	"net/http"
)

// NamedTagListController ...
type NamedTagListController interface {
	GetNamedTagLists() http.Handler
	CreateNamedTagList() http.Handler
	DeleteNamedTagLists() http.Handler
}

type namedTagListController struct {
	namedTagListRepository NamedTagListRepository
}

// NewNamedTagListController ...
func NewNamedTagListController(namedTagListRepository NamedTagListRepository) NamedTagListController {
	return &namedTagListController{
		namedTagListRepository,
	}
}

func (c *namedTagListController) GetNamedTagLists() http.Handler {
	return http.HandlerFunc(
		func(rw http.ResponseWriter, r *http.Request) {
			bytes, err := json.Marshal(
				c.namedTagListRepository.FindAll(),
			)
			if err != nil {
				panic(err)
			}
			rw.Write(bytes)
		},
	)
}

func (c *namedTagListController) CreateNamedTagList() http.Handler {
	return http.HandlerFunc(
		func(rw http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			var namedTagList NamedTagList
			if err := json.NewDecoder(r.Body).Decode(&namedTagList); err != nil {
				panic(err)
			}
			if c.namedTagListRepository.Create(namedTagList) != nil {
				rw.WriteHeader(500)
			} else {
				rw.WriteHeader(201)
			}
		},
	)
}

func (c *namedTagListController) DeleteNamedTagLists() http.Handler {
	return http.HandlerFunc(
		func(rw http.ResponseWriter, r *http.Request) {
			err := c.namedTagListRepository.DeleteAll()
			if err != nil {
				rw.WriteHeader(500)
			} else {
				rw.WriteHeader(204)
			}
		},
	)
}
