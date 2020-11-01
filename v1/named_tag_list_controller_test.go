package v1

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
)

type stubNamedTagListRepository struct {
	dummyNamedTagList NamedTagList
	willError         bool
}

func (r *stubNamedTagListRepository) FindAll() ([]NamedTagList, error) {
	if r.willError {
		return nil, errors.New("there was an error")
	}
	return []NamedTagList{r.dummyNamedTagList}, nil
}

func (r *stubNamedTagListRepository) Create(namedTagList NamedTagList) error {
	if r.willError && reflect.DeepEqual(namedTagList, r.dummyNamedTagList) {
		return errors.New("there was an error")
	}
	return nil
}

func (r *stubNamedTagListRepository) DeleteAll() error {
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

	t.Run("GET when repository has error", func(t *testing.T) {
		controller := NewNamedTagListController(
			&stubNamedTagListRepository{
				willError: true,
			},
		)

		request, _ := http.NewRequest(http.MethodGet, "/", nil)
		response := httptest.NewRecorder()
		controller.GetNamedTagLists().ServeHTTP(response, request)

		gotStatusCode := response.Result().StatusCode
		wantStatusCode := 500

		if gotStatusCode != wantStatusCode {
			t.Errorf("got status code %d want %d", gotStatusCode, wantStatusCode)
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
	})

	t.Run("POST when request body malformed", func(t *testing.T) {
		controller := NewNamedTagListController(
			&stubNamedTagListRepository{},
		)

		request, _ := http.NewRequest(http.MethodPost, "/", strings.NewReader("{\"garbalooy\":\"gook"))
		response := httptest.NewRecorder()
		controller.CreateNamedTagList().ServeHTTP(response, request)

		gotStatusCode := response.Result().StatusCode
		wantStatusCode := 400

		if gotStatusCode != wantStatusCode {
			t.Errorf("got status code %d want %d", gotStatusCode, wantStatusCode)
		}
	})

	t.Run("POST when repository has error", func(t *testing.T) {
		controller := NewNamedTagListController(
			&stubNamedTagListRepository{
				dummyNamedTagList: dummyNamedTagList,
				willError:         true,
			},
		)

		requestBody, err := json.Marshal(dummyNamedTagList)
		if err != nil {
			t.Fatal(err)
		}

		request, _ := http.NewRequest(http.MethodPost, "/", bytes.NewBuffer(requestBody))
		response := httptest.NewRecorder()
		controller.CreateNamedTagList().ServeHTTP(response, request)

		gotStatusCode := response.Result().StatusCode
		wantStatusCode := 500

		if gotStatusCode != wantStatusCode {
			t.Errorf("got status code %d want %d", gotStatusCode, wantStatusCode)
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
