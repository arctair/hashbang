package v1

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

type stubNamedTagListController struct {
}

func (c *stubNamedTagListController) GetNamedTagLists() http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("the named tag list controller body / get method"))
		},
	)
}

func (c *stubNamedTagListController) CreateNamedTagList() http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("the named tag list controller body / post method"))
		},
	)
}

func (c *stubNamedTagListController) DeleteNamedTagLists() http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("the named tag list controller body / delete method"))
		},
	)
}

type stubVersionController struct {
}

func (c *stubVersionController) HandlerFunc() http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("the version controller body"))
		},
	)
}

func TestRouter(t *testing.T) {
	router := NewRouter(
		&stubNamedTagListController{},
		&stubVersionController{},
	)

	t.Run("Route not found", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, "/", nil)
		response := httptest.NewRecorder()

		router.ServeHTTP(response, request)

		gotStatusCode := response.Result().StatusCode
		wantStatusCode := 404

		if gotStatusCode != wantStatusCode {
			t.Errorf("got status code %d want %d", gotStatusCode, wantStatusCode)
		}

		gotBody := string(response.Body.Bytes())
		wantBody := "404 page not found\n"

		if gotBody != wantBody {
			t.Errorf("got body %s want %s", gotBody, wantBody)
		}
	})

	t.Run("Route GET /namedTagLists to named tag list controller", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, "/namedTagLists", nil)
		response := httptest.NewRecorder()

		router.ServeHTTP(response, request)

		gotStatusCode := response.Result().StatusCode
		wantStatusCode := 200

		if gotStatusCode != wantStatusCode {
			t.Errorf("got status code %d want %d", gotStatusCode, wantStatusCode)
		}

		gotBody := string(response.Body.Bytes())
		wantBody := "the named tag list controller body / get method"

		if gotBody != wantBody {
			t.Errorf("got body %s want %s", gotBody, wantBody)
		}
	})

	t.Run("Route POST /namedTagLists to named tag list controller", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodPost, "/namedTagLists", nil)
		response := httptest.NewRecorder()

		router.ServeHTTP(response, request)

		gotStatusCode := response.Result().StatusCode
		wantStatusCode := 200

		if gotStatusCode != wantStatusCode {
			t.Errorf("got status code %d want %d", gotStatusCode, wantStatusCode)
		}

		gotBody := string(response.Body.Bytes())
		wantBody := "the named tag list controller body / post method"

		if gotBody != wantBody {
			t.Errorf("got body %s want %s", gotBody, wantBody)
		}
	})

	t.Run("Route DELETE /namedTagLists to named tag list controller", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodDelete, "/namedTagLists", nil)
		response := httptest.NewRecorder()

		router.ServeHTTP(response, request)

		gotStatusCode := response.Result().StatusCode
		wantStatusCode := 200

		if gotStatusCode != wantStatusCode {
			t.Errorf("got status code %d want %d", gotStatusCode, wantStatusCode)
		}

		gotBody := string(response.Body.Bytes())
		wantBody := "the named tag list controller body / delete method"

		if gotBody != wantBody {
			t.Errorf("got body %s want %s", gotBody, wantBody)
		}
	})

	t.Run("Route /version to version controller", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, "/version", nil)
		response := httptest.NewRecorder()

		router.ServeHTTP(response, request)

		gotStatusCode := response.Result().StatusCode
		wantStatusCode := 200

		if gotStatusCode != wantStatusCode {
			t.Errorf("got status code %d want %d", gotStatusCode, wantStatusCode)
		}

		gotBody := string(response.Body.Bytes())
		wantBody := "the version controller body"

		if gotBody != wantBody {
			t.Errorf("got body %s want %s", gotBody, wantBody)
		}
	})
}
