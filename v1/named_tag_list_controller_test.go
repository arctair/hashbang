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

type stubLogger struct {
	errors []string
}

func stubLoggerNew() *stubLogger {
	return &stubLogger{
		errors: []string{},
	}
}

func (l *stubLogger) Error(err error) {
	l.errors = append(l.errors, fmt.Sprint(err))
}

type stubNamedTagListRepositoryForController struct {
	NamedTagListRepository

	withBuckets      []string
	withIds          []string
	withNamedTagList NamedTagList
	willError        string

	err error
}

func (r *stubNamedTagListRepositoryForController) FindAll(buckets []string) ([]NamedTagList, error) {
	requestMatched := reflect.DeepEqual(buckets, r.withBuckets)
	if !requestMatched {
		r.err = fmt.Errorf("Stub got buckets %v want %v", buckets, r.withBuckets)
	}
	if requestMatched == (r.willError == "FindAll") {
		return nil, errors.New("there was an error")
	}
	return []NamedTagList{r.withNamedTagList}, nil
}

func (r *stubNamedTagListRepositoryForController) FindAllOld() ([]NamedTagList, error) {
	return nil, errors.New("do not call")
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
	withBucket       string
	withNamedTagList NamedTagList
	willError        string

	err error
}

func (r *stubNamedTagListService) CreateOld(ntl NamedTagList) (*NamedTagList, error) {
	return nil, errors.New("do not call")
}

func (r *stubNamedTagListService) Create(bucket string, ntl NamedTagList) (*NamedTagList, error) {
	requestMatched := bucket == r.withBucket && reflect.DeepEqual(ntl, r.withNamedTagList)
	if !requestMatched {
		r.err = fmt.Errorf("Stub got bucket %s want %v got ntl %+v want %+v", bucket, r.withBucket, ntl, r.withNamedTagList)
	}
	if requestMatched == (r.willError == "Create") {
		return nil, errors.New("there was an error")
	}
	return &r.withNamedTagList, nil
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
		repository := &stubNamedTagListRepositoryForController{
			withBuckets:      []string{"red", "blue"},
			withNamedTagList: dummyNamedTagList,
		}
		controller := NewNamedTagListController(
			stubLoggerNew(),
			repository,
			&stubNamedTagListService{},
		)

		request, _ := http.NewRequest(http.MethodGet, "/?bucket=red&bucket=blue", nil)
		response := httptest.NewRecorder()
		controller.GetNamedTagLists().ServeHTTP(response, request)

		if repository.err != nil {
			t.Error(repository.err)
		}

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
		logger := stubLoggerNew()
		repository := &stubNamedTagListRepositoryForController{
			withBuckets: []string{"bucket"},
			willError:   "FindAll",
		}
		controller := NewNamedTagListController(
			logger,
			repository,
			&stubNamedTagListService{},
		)

		request, _ := http.NewRequest(http.MethodGet, "/?bucket=bucket", nil)
		response := httptest.NewRecorder()
		controller.GetNamedTagLists().ServeHTTP(response, request)

		if repository.err != nil {
			t.Error(repository.err)
		}

		gotStatusCode := response.Result().StatusCode
		wantStatusCode := 500

		if gotStatusCode != wantStatusCode {
			t.Errorf("got status code %d want %d", gotStatusCode, wantStatusCode)
		}

		wantErrorf := []string{"there was an error"}
		if !reflect.DeepEqual(logger.errors, wantErrorf) {
			t.Errorf("got logger.Errorf %+v want %+v", logger.errors, wantErrorf)
		}
	})

	t.Run("GET with no buckets", func(t *testing.T) {
		controller := NewNamedTagListController(
			stubLoggerNew(),
			&stubNamedTagListRepositoryForController{},
			&stubNamedTagListService{},
		)

		request, _ := http.NewRequest(http.MethodGet, "/", nil)
		response := httptest.NewRecorder()
		controller.GetNamedTagLists().ServeHTTP(response, request)

		gotStatusCode := response.Result().StatusCode
		wantStatusCode := 400

		if gotStatusCode != wantStatusCode {
			t.Errorf("got status code %d want %d", gotStatusCode, wantStatusCode)
		}

		var gotResponseBody map[string]string
		if err := json.NewDecoder(response.Body).Decode(&gotResponseBody); err != nil {
			t.Fatal(err)
		}

		wantResponseBody := map[string]string{"error": "bucket query parameter is required"}

		if !reflect.DeepEqual(gotResponseBody, wantResponseBody) {
			t.Errorf("got response body %+v want %+v", gotResponseBody, wantResponseBody)
		}
	})
	t.Run("POST", func(t *testing.T) {
		service := &stubNamedTagListService{
			withBucket:       "bucket",
			withNamedTagList: dummyNamedTagList,
		}
		controller := NewNamedTagListController(
			stubLoggerNew(),
			&stubNamedTagListRepositoryForController{},
			service,
		)

		requestBody, err := json.Marshal(dummyNamedTagList)
		if err != nil {
			t.Fatal(err)
		}

		request, _ := http.NewRequest(http.MethodPost, "/?bucket=bucket", bytes.NewBuffer(requestBody))
		response := httptest.NewRecorder()
		controller.CreateNamedTagList().ServeHTTP(response, request)

		if service.err != nil {
			t.Error(service.err)
		}

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
			stubLoggerNew(),
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
		logger := stubLoggerNew()
		service := &stubNamedTagListService{
			withBucket:       "bucket",
			withNamedTagList: dummyNamedTagList,
			willError:        "Create",
		}
		controller := NewNamedTagListController(
			logger,
			&stubNamedTagListRepositoryForController{},
			service,
		)

		requestBody, err := json.Marshal(dummyNamedTagList)
		if err != nil {
			t.Fatal(err)
		}

		request, _ := http.NewRequest(http.MethodPost, "/?bucket=bucket", bytes.NewBuffer(requestBody))
		response := httptest.NewRecorder()
		controller.CreateNamedTagList().ServeHTTP(response, request)

		if service.err != nil {
			t.Error(service.err)
		}

		gotStatusCode := response.Result().StatusCode
		wantStatusCode := 500

		if gotStatusCode != wantStatusCode {
			t.Errorf("got status code %d want %d", gotStatusCode, wantStatusCode)
		}

		wantErrorf := []string{"there was an error"}
		if !reflect.DeepEqual(logger.errors, wantErrorf) {
			t.Errorf("got logger.Errorf %+v want %+v", logger.errors, wantErrorf)
		}
	})

	t.Run("POST when bucket is empty", func(t *testing.T) {
		controller := NewNamedTagListController(
			stubLoggerNew(),
			&stubNamedTagListRepositoryForController{},
			&stubNamedTagListService{},
		)

		requestBody, err := json.Marshal(dummyNamedTagList)
		if err != nil {
			t.Fatal(err)
		}

		request, _ := http.NewRequest(http.MethodPost, "/", bytes.NewBuffer(requestBody))
		response := httptest.NewRecorder()
		controller.CreateNamedTagList().ServeHTTP(response, request)

		gotStatusCode := response.Result().StatusCode
		wantStatusCode := 400

		if gotStatusCode != wantStatusCode {
			t.Errorf("got status code %d want %d", gotStatusCode, wantStatusCode)
		}

		var gotResponseBody map[string]string
		if err := json.NewDecoder(response.Body).Decode(&gotResponseBody); err != nil {
			t.Fatal(err)
		}

		wantResponseBody := map[string]string{"error": "bucket query parameter is required"}

		if !reflect.DeepEqual(gotResponseBody, wantResponseBody) {
			t.Errorf("got response body %+v want %+v", gotResponseBody, wantResponseBody)
		}
	})

	t.Run("POST when more than one bucket", func(t *testing.T) {
		controller := NewNamedTagListController(
			stubLoggerNew(),
			&stubNamedTagListRepositoryForController{},
			&stubNamedTagListService{},
		)

		requestBody, err := json.Marshal(dummyNamedTagList)
		if err != nil {
			t.Fatal(err)
		}

		request, _ := http.NewRequest(http.MethodPost, "/?bucket=one&bucket=two", bytes.NewBuffer(requestBody))
		response := httptest.NewRecorder()
		controller.CreateNamedTagList().ServeHTTP(response, request)

		gotStatusCode := response.Result().StatusCode
		wantStatusCode := 400

		if gotStatusCode != wantStatusCode {
			t.Errorf("got status code %d want %d", gotStatusCode, wantStatusCode)
		}

		var gotResponseBody map[string]string
		if err := json.NewDecoder(response.Body).Decode(&gotResponseBody); err != nil {
			t.Fatal(err)
		}

		wantResponseBody := map[string]string{"error": "no more than one bucket must be supplied"}

		if !reflect.DeepEqual(gotResponseBody, wantResponseBody) {
			t.Errorf("got named tag list %+v want %+v", gotResponseBody, wantResponseBody)
		}
	})

	t.Run("PUT", func(t *testing.T) {
		stubRepository := &stubNamedTagListRepositoryForController{
			withIds:          []string{"deadbeef-dead-beef-dead-beefdeadbeef"},
			withNamedTagList: dummyNamedTagList,
		}
		controller := NewNamedTagListController(
			stubLoggerNew(),
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
			stubLoggerNew(),
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
		logger := stubLoggerNew()
		controller := NewNamedTagListController(
			logger,
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

		wantErrorf := []string{"there was an error"}
		if !reflect.DeepEqual(logger.errors, wantErrorf) {
			t.Errorf("got logger.Errorf %+v want %+v", logger.errors, wantErrorf)
		}
	})

	t.Run("DELETE by ids", func(t *testing.T) {
		controller := NewNamedTagListController(
			stubLoggerNew(),
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
		logger := stubLoggerNew()
		controller := NewNamedTagListController(
			logger,
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

		wantErrorf := []string{"there was an error"}
		if !reflect.DeepEqual(logger.errors, wantErrorf) {
			t.Errorf("got logger.Errorf %+v want %+v", logger.errors, wantErrorf)
		}
	})

	t.Run("DELETE all", func(t *testing.T) {
		repository := &stubNamedTagListRepositoryForController{}
		controller := NewNamedTagListController(
			stubLoggerNew(),
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
		logger := stubLoggerNew()
		controller := NewNamedTagListController(
			logger,
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

		wantErrorf := []string{"there was an error"}
		if !reflect.DeepEqual(logger.errors, wantErrorf) {
			t.Errorf("got logger.Errorf %+v want %+v", logger.errors, wantErrorf)
		}
	})
}
