package v1

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

type stubNamedTagListRepository struct {
	created    NamedTagList
	deletedAll bool
}

func (r *stubNamedTagListRepository) FindAll() []NamedTagList {
	return []NamedTagList{
		{
			Name: "tag list name",
			Tags: []string{
				"#windy",
				"#tdd",
			},
		},
	}
}

func (r *stubNamedTagListRepository) Create(namedTagList NamedTagList) {
	r.created = namedTagList
}

func (r *stubNamedTagListRepository) DeleteAll() {
	r.deletedAll = true
}

func TestNamedTagListController(t *testing.T) {
	t.Run("GET", func(t *testing.T) {
		controller := NewNamedTagListController(
			&stubNamedTagListRepository{},
		)

		request, _ := http.NewRequest(http.MethodGet, "/", nil)
		response := httptest.NewRecorder()
		controller.GetNamedTagLists().ServeHTTP(response, request)

		gotStatusCode := response.Result().StatusCode
		wantStatusCode := 200

		if gotStatusCode != wantStatusCode {
			t.Errorf("got status code %d want %d", gotStatusCode, wantStatusCode)
		}

		var gotNamedTagLists []NamedTagList
		if err := json.NewDecoder(response.Body).Decode(&gotNamedTagLists); err != nil {
			t.Fatal(err)
		}

		wantNamedTagLists := []NamedTagList{
			{
				Name: "tag list name",
				Tags: []string{
					"#windy",
					"#tdd",
				},
			},
		}

		if !reflect.DeepEqual(gotNamedTagLists, wantNamedTagLists) {
			t.Errorf("got named tag lists %q want %q", gotNamedTagLists, wantNamedTagLists)
		}
	})

	t.Run("POST", func(t *testing.T) {
		repository := &stubNamedTagListRepository{}
		controller := NewNamedTagListController(
			repository,
		)

		namedTagList := NamedTagList{
			Name: "tag list name",
			Tags: []string{
				"#windy",
				"#tdd",
			},
		}
		requestBody, err := json.Marshal(namedTagList)
		if err != nil {
			t.Fatal(err)
		}

		request, _ := http.NewRequest(http.MethodPost, "/", bytes.NewBuffer(requestBody))
		response := httptest.NewRecorder()
		controller.CreateNamedTagList().ServeHTTP(response, request)

		gotStatusCode := response.Result().StatusCode
		wantStatusCode := 201

		if gotStatusCode != wantStatusCode {
			t.Errorf("got status code %d want %d", gotStatusCode, wantStatusCode)
		}

		gotCreated := repository.created
		wantCreated := NamedTagList{
			Name: "tag list name",
			Tags: []string{
				"#windy",
				"#tdd",
			},
		}

		if !reflect.DeepEqual(gotCreated, wantCreated) {
			t.Errorf("got created %q want %q", gotCreated, wantCreated)
		}
	})

	t.Run("DELETE", func(t *testing.T) {
		repository := &stubNamedTagListRepository{}
		controller := NewNamedTagListController(
			repository,
		)

		request, _ := http.NewRequest(http.MethodDelete, "/namedTagLists", nil)
		response := httptest.NewRecorder()
		controller.DeleteNamedTagLists().ServeHTTP(response, request)

		gotStatusCode := response.Result().StatusCode
		wantStatusCode := 204

		if gotStatusCode != wantStatusCode {
			t.Errorf("got status code %d want %d", gotStatusCode, wantStatusCode)
		}

		gotDeletedAll := repository.deletedAll

		if !gotDeletedAll {
			t.Errorf("got deleted all false wanted true")
		}
	})
}
