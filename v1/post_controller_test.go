package v1

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestPostController(t *testing.T) {
	t.Run("GET returns posts", func(t *testing.T) {
		postController := NewPostController()

		request, _ := http.NewRequest(http.MethodGet, "/", nil)
		response := httptest.NewRecorder()
		postController.HandlerFunc().ServeHTTP(response, request)

		gotStatusCode := response.Result().StatusCode
		wantStatusCode := 200

		if gotStatusCode != wantStatusCode {
			t.Errorf("got status code %d want %d", gotStatusCode, wantStatusCode)
		}

		var gotBody []Post
		if err := json.NewDecoder(response.Body).Decode(&gotBody); err != nil {
			t.Fatal(err)
		}

		wantBody := []Post{
			{
				ImageUri: "https://images.unsplash.com/photo-1603316851229-26637b4bd1b8?ixlib=rb-1.2.1&ixid=eyJhcHBfaWQiOjEyMDd9&auto=format&fit=crop&w=1400&q=80",
				Tags: []string{
					"#windy",
					"#tdd",
				},
			},
		}

		if !reflect.DeepEqual(gotBody, wantBody) {
			t.Errorf("got body %q want %q", gotBody, wantBody)
		}
	})
}
