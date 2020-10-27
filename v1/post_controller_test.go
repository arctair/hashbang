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
	created    Post
	deletedAll bool
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

func (r *stubPostRepository) Create(post Post) Post {
	post.Id = "80000000-0000-0000-0000-000000000000"
	return post
}

func (r *stubPostRepository) DeleteAll() {
	r.deletedAll = true
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

		var gotPosts []Post
		if err := json.NewDecoder(response.Body).Decode(&gotPosts); err != nil {
			t.Fatal(err)
		}

		wantPosts := []Post{
			{
				ImageUri: "https://images.unsplash.com/photo-1603316851229-26637b4bd1b8?ixlib=rb-1.2.1&ixid=eyJhcHBfaWQiOjEyMDd9&auto=format&fit=crop&w=1400&q=80",
				Tags: []string{
					"#windy",
					"#tdd",
				},
			},
		}

		if !reflect.DeepEqual(gotPosts, wantPosts) {
			t.Errorf("got posts %q want %q", gotPosts, wantPosts)
		}
	})

	t.Run("POST creates post", func(t *testing.T) {
		repository := &stubPostRepository{}
		postController := NewPostController(
			repository,
		)

		post := Post{
			ImageUri: "https://images.unsplash.com/photo-1603316851229-26637b4bd1b8?ixlib=rb-1.2.1&ixid=eyJhcHBfaWQiOjEyMDd9&auto=format&fit=crop&w=1400&q=80",
			Tags: []string{
				"#windy",
				"#tdd",
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

		var gotPost Post
		if err := json.NewDecoder(response.Body).Decode(&gotPost); err != nil {
			t.Fatal(err)
		}

		wantPost := Post{
			Id:       "80000000-0000-0000-0000-000000000000",
			ImageUri: "https://images.unsplash.com/photo-1603316851229-26637b4bd1b8?ixlib=rb-1.2.1&ixid=eyJhcHBfaWQiOjEyMDd9&auto=format&fit=crop&w=1400&q=80",
			Tags: []string{
				"#windy",
				"#tdd",
			},
		}

		if !reflect.DeepEqual(gotPost, wantPost) {
			t.Errorf("got post %+q want %+q", gotPost, wantPost)
		}
	})

	t.Run("DELETE deletes posts", func(t *testing.T) {
		repository := &stubPostRepository{}
		postController := NewPostController(
			repository,
		)

		request, _ := http.NewRequest(http.MethodDelete, "/posts?id=8c907ab9-fef8-43ab-9103-b19aabfb40b2", nil)
		response := httptest.NewRecorder()
		postController.DeletePost().ServeHTTP(response, request)

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
