package v1

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
)

type stubNamedTagListRepositoryForController struct {
	NamedTagListRepository

	withIds          []string
	withNamedTagList NamedTagList
	willError        string

	err error
}

func (r *stubNamedTagListRepositoryForController) FindAll() ([]NamedTagList, error) {
	if r.willError == "FindAll" {
		return nil, errors.New("there was an error")
	}
	return []NamedTagList{r.withNamedTagList}, nil
}

func (r *stubNamedTagListRepositoryForController) ReplaceByIds(ids []string, ntl NamedTagList) error {
	requestMatched := reflect.DeepEqual(ids, r.withIds) && reflect.DeepEqual(ntl, r.withNamedTagList)
	willError := (r.willError == "ReplaceByIds")
	if !requestMatched {
		r.err = fmt.Errorf("Stub got ids %v want %v got ntl %+v want %+v", ids, r.withIds, ntl, r.withNamedTagList)
	}
	if requestMatched == willError {
		return fmt.Errorf("there was an error")
	}
	return nil
}

func (r *stubNamedTagListRepositoryForController) DeleteByIds(ids []string) error {
	if r.willError == "DeleteByIds" {
		return errors.New("there was an error")
	}
	return nil
}

func (r *stubNamedTagListRepositoryForController) DeleteAll() error {
	if r.willError == "DeleteAll" {
		return errors.New("there was an error")
	}
	return nil
}

type stubNamedTagListService struct {
	dummyNamedTagList NamedTagList
	willError         string
}

func (r *stubNamedTagListService) Create(namedTagList NamedTagList) (*NamedTagList, error) {
	if r.willError == "Create" && reflect.DeepEqual(namedTagList, r.dummyNamedTagList) {
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
				withNamedTagList: dummyNamedTagList,
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
				willError: "FindAll",
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
				willError:         "Create",
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

	t.Run("PUT", func(t *testing.T) {
		stubRepository := &stubNamedTagListRepositoryForController{
			withIds:          []string{"deadbeef-dead-beef-dead-beefdeadbeef"},
			withNamedTagList: dummyNamedTagList,
		}
		controller := NewNamedTagListController(
			stubRepository,
			&stubNamedTagListService{},
		)

		requestBody, err := json.Marshal(dummyNamedTagList)
		if err != nil {
			t.Fatal(err)
		}

		request, _ := http.NewRequest(
			http.MethodPut,
			"/?id=deadbeef-dead-beef-dead-beefdeadbeef",
			bytes.NewBuffer(requestBody),
		)
		response := httptest.NewRecorder()
		controller.ReplaceNamedTagLists().ServeHTTP(response, request)

		gotStatusCode := response.Result().StatusCode
		wantStatusCode := 204

		if gotStatusCode != wantStatusCode {
			t.Errorf("got status code %d want %d", gotStatusCode, wantStatusCode)
		}

		if stubRepository.err != nil {
			t.Error(stubRepository.err)
		}
	})

	t.Run("PUT when request body malformed", func(t *testing.T) {
		controller := NewNamedTagListController(
			&stubNamedTagListRepositoryForController{},
			&stubNamedTagListService{},
		)

		request, _ := http.NewRequest(http.MethodPut, "/", strings.NewReader("{\"garbalooy\":\"gook"))
		response := httptest.NewRecorder()
		controller.ReplaceNamedTagLists().ServeHTTP(response, request)

		gotStatusCode := response.Result().StatusCode
		wantStatusCode := 400

		if gotStatusCode != wantStatusCode {
			t.Errorf("got status code %d want %d", gotStatusCode, wantStatusCode)
		}
	})

	t.Run("PUT when repository has error", func(t *testing.T) {
		stubRepository := &stubNamedTagListRepositoryForController{
			withIds:          []string{"deadbeef-dead-beef-dead-beefdeadbeef"},
			withNamedTagList: dummyNamedTagList,
			willError:        "ReplaceByIds",
		}
		controller := NewNamedTagListController(
			stubRepository,
			&stubNamedTagListService{},
		)

		requestBody, err := json.Marshal(dummyNamedTagList)
		if err != nil {
			t.Fatal(err)
		}

		request, _ := http.NewRequest(
			http.MethodPut,
			"/?id=deadbeef-dead-beef-dead-beefdeadbeef",
			bytes.NewBuffer(requestBody),
		)
		response := httptest.NewRecorder()
		controller.ReplaceNamedTagLists().ServeHTTP(response, request)

		gotStatusCode := response.Result().StatusCode
		wantStatusCode := 500

		if gotStatusCode != wantStatusCode {
			t.Errorf("got status code %d want %d", gotStatusCode, wantStatusCode)
		}

		if stubRepository.err != nil {
			t.Error(stubRepository.err)
		}
	})

	t.Run("DELETE by ids", func(t *testing.T) {
		controller := NewNamedTagListController(
			&stubNamedTagListRepositoryForController{},
			&stubNamedTagListService{},
		)

		request, _ := http.NewRequest(http.MethodDelete, "/namedTagLists?id=0b491dfc-3969-4ae3-83dd-83fae3b0f56e", nil)
		response := httptest.NewRecorder()
		controller.DeleteNamedTagLists().ServeHTTP(response, request)

		gotStatusCode := response.Result().StatusCode
		wantStatusCode := 204

		if gotStatusCode != wantStatusCode {
			t.Errorf("got status code %d want %d", gotStatusCode, wantStatusCode)
		}
	})

	t.Run("DELETE by ids when repository has error", func(t *testing.T) {
		controller := NewNamedTagListController(
			&stubNamedTagListRepositoryForController{
				willError: "DeleteByIds",
			},
			&stubNamedTagListService{},
		)

		request, _ := http.NewRequest(http.MethodDelete, "/namedTagLists?id=0b491dfc-3969-4ae3-83dd-83fae3b0f56e", nil)
		response := httptest.NewRecorder()
		controller.DeleteNamedTagLists().ServeHTTP(response, request)

		gotStatusCode := response.Result().StatusCode
		wantStatusCode := 500

		if gotStatusCode != wantStatusCode {
			t.Errorf("got status code %d want %d", gotStatusCode, wantStatusCode)
		}
	})

	t.Run("DELETE all", func(t *testing.T) {
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

	t.Run("DELETE all when repository has error", func(t *testing.T) {
		controller := NewNamedTagListController(
			&stubNamedTagListRepositoryForController{
				willError: "DeleteAll",
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
