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

type stubNamedTagListRepositoryForController struct {
	NamedTagListRepository

	dummyNamedTagList NamedTagList
	willError         bool
}

func (r *stubNamedTagListRepositoryForController) FindAll() ([]NamedTagList, error) {
	if r.willError {
		return nil, errors.New("there was an error")
	}
	return []NamedTagList{r.dummyNamedTagList}, nil
}

func (r *stubNamedTagListRepositoryForController) DeleteAll() error {
	if r.willError {
		return errors.New("there was an error")
	}
	return nil
}

type stubNamedTagListService struct {
	dummyNamedTagList NamedTagList
	willError         bool
}

func (r *stubNamedTagListService) Create(namedTagList NamedTagList) (*NamedTagList, error) {
	if r.willError && reflect.DeepEqual(namedTagList, r.dummyNamedTagList) {
		return nil, errors.New("there was an error")
	}
	return &r.dummyNamedTagList, nil
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
			&stubNamedTagListRepositoryForController{
				dummyNamedTagList: dummyNamedTagList,
			},
			&stubNamedTagListService{},
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
			&stubNamedTagListRepositoryForController{
				willError: true,
			},
			&stubNamedTagListService{},
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
		controller := NewNamedTagListController(
			&stubNamedTagListRepositoryForController{},
			&stubNamedTagListService{
				dummyNamedTagList: dummyNamedTagList,
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
		wantStatusCode := 201

		if gotStatusCode != wantStatusCode {
			t.Errorf("got status code %d want %d", gotStatusCode, wantStatusCode)
		}

		var gotNamedTagList NamedTagList
		if err := json.NewDecoder(response.Body).Decode(&gotNamedTagList); err != nil {
			t.Fatal(err)
		}

		wantNamedTagList := dummyNamedTagList

		if !reflect.DeepEqual(gotNamedTagList, wantNamedTagList) {
			t.Errorf("got named tag list %+v want %+v", gotNamedTagList, wantNamedTagList)
		}
	})

	t.Run("POST when request body malformed", func(t *testing.T) {
		controller := NewNamedTagListController(
			&stubNamedTagListRepositoryForController{},
			&stubNamedTagListService{},
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
			&stubNamedTagListRepositoryForController{},
			&stubNamedTagListService{
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
		repository := &stubNamedTagListRepositoryForController{}
		controller := NewNamedTagListController(
			repository,
			&stubNamedTagListService{},
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
			&stubNamedTagListRepositoryForController{
				willError: true,
			},
			&stubNamedTagListService{},
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
