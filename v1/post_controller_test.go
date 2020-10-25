package v1

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

type stubPostRepository struct {
	created Post
}

func (r *stubPostRepository) FindAll() []Post {
	return []Post{
		{
			ImageUri: "https://images.unsplash.com/photo-1603316851229-26637b4bd1b8?ixlib=rb-1.2.1&ixid=eyJhcHBfaWQiOjEyMDd9&auto=format&fit=crop&w=1400&q=80",
			Tags: []string{
				"#windy",
				"#tdd",
			},
		},
	}
}

func (r *stubPostRepository) Create(post Post) {
	r.created = post
}

func TestPostController(t *testing.T) {
	t.Run("GET returns posts", func(t *testing.T) {
		postController := NewPostController(
			&stubPostRepository{},
		)

		request, _ := http.NewRequest(http.MethodGet, "/", nil)
		response := httptest.NewRecorder()
		postController.GetPosts().ServeHTTP(response, request)

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

	t.Run("POST creates post", func(t *testing.T) {
		repository := &stubPostRepository{}
		postController := NewPostController(
			repository,
		)

		post := []Post{
			{
				ImageUri: "https://images.unsplash.com/photo-1603316851229-26637b4bd1b8?ixlib=rb-1.2.1&ixid=eyJhcHBfaWQiOjEyMDd9&auto=format&fit=crop&w=1400&q=80",
				Tags: []string{
					"#windy",
					"#tdd",
				},
			},
		}
		requestBody, err := json.Marshal(post)
		if err != nil {
			t.Fatal(err)
		}

		request, _ := http.NewRequest(http.MethodPost, "/", bytes.NewBuffer(requestBody))
		response := httptest.NewRecorder()
		postController.CreatePost().ServeHTTP(response, request)

		gotStatusCode := response.Result().StatusCode
		wantStatusCode := 201

		if gotStatusCode != wantStatusCode {
			t.Errorf("got status code %d want %d", gotStatusCode, wantStatusCode)
		}

		gotCreated := repository.created
		wantCreated := Post{
			ImageUri: "https://images.unsplash.com/photo-1603316851229-26637b4bd1b8?ixlib=rb-1.2.1&ixid=eyJhcHBfaWQiOjEyMDd9&auto=format&fit=crop&w=1400&q=80",
			Tags: []string{
				"#windy",
				"#tdd",
			},
		}

		if !reflect.DeepEqual(gotCreated, wantCreated) {
			t.Errorf("got created %q want %q", gotCreated, wantCreated)
		}
	})
}
