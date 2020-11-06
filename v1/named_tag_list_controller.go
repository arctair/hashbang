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
	namedTagListService    NamedTagListService
}

func (c *namedTagListController) GetNamedTagLists() http.Handler {
	return http.HandlerFunc(
		func(rw http.ResponseWriter, r *http.Request) {
			var (
				namedTagLists []NamedTagList
				err           error
			)
			if namedTagLists, err = c.namedTagListRepository.FindAll(); err != nil {
				rw.WriteHeader(500)
			}
			bytes, err := json.Marshal(namedTagLists)
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
			var (
				namedTagList *NamedTagList
				err          error
			)
			if json.NewDecoder(r.Body).Decode(&namedTagList) != nil {
				rw.WriteHeader(http.StatusBadRequest)
			} else if namedTagList, err = c.namedTagListService.Create(*namedTagList); err != nil {
				rw.WriteHeader(http.StatusInternalServerError)
			} else {
				rw.WriteHeader(http.StatusCreated)
				json.NewEncoder(rw).Encode(namedTagList)
			}
		},
	)
}

func (c *namedTagListController) DeleteNamedTagLists() http.Handler {
	return http.HandlerFunc(
		func(rw http.ResponseWriter, r *http.Request) {
			var err error
			ids := r.URL.Query()["id"]
			if len(ids) > 0 {
				err = c.namedTagListRepository.DeleteByIds(ids)
			} else {
				err = c.namedTagListRepository.DeleteAll()
			}
			if err != nil {
				rw.WriteHeader(500)
			} else {
				rw.WriteHeader(204)
			}
		},
	)
}

// NewNamedTagListController ...
func NewNamedTagListController(
	namedTagListRepository NamedTagListRepository,
	namedTagListService NamedTagListService,
) NamedTagListController {
	return &namedTagListController{
		namedTagListRepository,
		namedTagListService,
	}
}
