package v1

import (
	"encoding/json"
	"net/http"
)

// NamedTagListController ...
type NamedTagListController interface {
	GetNamedTagLists() http.Handler
	CreateNamedTagList() http.Handler
	ReplaceNamedTagLists() http.Handler
	DeleteNamedTagLists() http.Handler
}

type namedTagListController struct {
	logger                 Logger
	namedTagListRepository NamedTagListRepository
	namedTagListService    NamedTagListService
}

func (c *namedTagListController) GetNamedTagLists() http.Handler {
	return http.HandlerFunc(
		func(rw http.ResponseWriter, r *http.Request) {
			buckets := r.URL.Query()["bucket"]
			if len(buckets) < 1 {
				rw.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(rw).Encode(map[string]string{"error": "bucket query parameter is required"})
				return
			}

			var (
				namedTagLists []NamedTagList
				err           error
			)
			if namedTagLists, err = c.namedTagListRepository.FindAll(buckets); err != nil {
				rw.WriteHeader(http.StatusInternalServerError)
				c.logger.Error(err)
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
			buckets := r.URL.Query()["bucket"]
			if len(buckets) < 1 {
				rw.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(rw).Encode(map[string]string{"error": "bucket query parameter is required"})
				return
			} else if len(buckets) > 1 {
				rw.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(rw).Encode(map[string]string{"error": "no more than one bucket must be supplied"})
				return
			}

			defer r.Body.Close()
			var (
				namedTagList *NamedTagList
				err          error
			)
			if json.NewDecoder(r.Body).Decode(&namedTagList) != nil {
				rw.WriteHeader(http.StatusBadRequest)
			} else if namedTagList, err = c.namedTagListService.Create(buckets[0], *namedTagList); err != nil {
				rw.WriteHeader(http.StatusInternalServerError)
				c.logger.Error(err)
			} else {
				rw.WriteHeader(http.StatusCreated)
				json.NewEncoder(rw).Encode(namedTagList)
			}
		},
	)
}

func (c *namedTagListController) ReplaceNamedTagLists() http.Handler {
	return http.HandlerFunc(
		func(rw http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			var namedTagList *NamedTagList
			if json.NewDecoder(r.Body).Decode(&namedTagList) != nil {
				rw.WriteHeader(http.StatusBadRequest)
			} else if err := c.namedTagListRepository.ReplaceByIds(r.URL.Query()["id"], *namedTagList); err != nil {
				rw.WriteHeader(500)
				c.logger.Error(err)
			} else {
				rw.WriteHeader(204)
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
				c.logger.Error(err)
			} else {
				rw.WriteHeader(204)
			}
		},
	)
}

// NewNamedTagListController ...
func NewNamedTagListController(
	logger Logger,
	namedTagListRepository NamedTagListRepository,
	namedTagListService NamedTagListService,
) NamedTagListController {
	return &namedTagListController{
		logger,
		namedTagListRepository,
		namedTagListService,
	}
}
