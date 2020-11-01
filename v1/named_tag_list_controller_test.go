package v1

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

type stubNamedTagListRepository struct {
	dummyNamedTagList NamedTagList
	created           NamedTagList
	deletedAll        bool
	willError         bool
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

func (r *stubNamedTagListRepository) DeleteAll() error {
	r.deletedAll = true
	if r.willError {
		return errors.New("there was an error")
	}
	return nil
}

func TestNamedTagListController(t *testing.T) {
	dummyNamedTagList := NamedTagList{
		Name: "tag list name",
		Tags: []string{
			"#windy",
			"#tdd",
		},
	}

	t.Run("GET", func(t *testing.T) {
		controller := NewNamedTagListController(
			&stubNamedTagListRepository{
				dummyNamedTagList: dummyNamedTagList,
			},
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
			dummyNamedTagList,
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

		requestBody, err := json.Marshal(dummyNamedTagList)
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
		wantCreated := dummyNamedTagList
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

	t.Run("DELETE when repository has error", func(t *testing.T) {
		controller := NewNamedTagListController(
			&stubNamedTagListRepository{
				willError: true,
			},
		)

		request, _ := http.NewRequest(http.MethodDelete, "/namedTagLists", nil)
		response := httptest.NewRecorder()
		controller.DeleteNamedTagLists().ServeHTTP(response, request)

		gotStatusCode := response.Result().StatusCode
		wantStatusCode := 500

		if gotStatusCode != wantStatusCode {
			t.Errorf("got status code %d want %d", gotStatusCode, wantStatusCode)
		}
	})
}
