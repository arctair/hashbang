package v1

import (
	"context"
	"reflect"
	"testing"

	"github.com/arctair/go-assertutil"
	"github.com/cockroachdb/cockroach-go/v2/testserver"
	"github.com/jackc/pgx/v4"
)

func TestPostRepository(t *testing.T) {
	testServer, err := testserver.NewTestServer()
	defer testServer.Stop()
	assertutil.NotError(t, err)

	connection, err := pgx.Connect(context.Background(), testServer.PGURL().String())
	assertutil.NotError(t, err)
	assertutil.NotError(t, Migrate(connection))

	postRepository := NewPostRepository(connection)

	t.Run("create, get, delete post", func(t *testing.T) {
		gotPosts := NewPostRepository(connection).FindAll()
		wantPosts := []Post{}

		if !reflect.DeepEqual(gotPosts, wantPosts) {
			t.Errorf("got %+v want %+v", gotPosts, wantPosts)
		}

		postRepository.Create(
			Post{
				ImageUri: "https://images.unsplash.com/photo-1603316851229-26637b4bd1b8?ixlib=rb-1.2.1&ixid=eyJhcHBfaWQiOjEyMDd9&auto=format&fit=crop&w=1400&q=80",
				Tags: []string{
					"#windy",
					"#tdd",
				},
			},
		)

		gotPosts = NewPostRepository(connection).FindAll()
		wantPosts = []Post{
			{
				ImageUri: "https://images.unsplash.com/photo-1603316851229-26637b4bd1b8?ixlib=rb-1.2.1&ixid=eyJhcHBfaWQiOjEyMDd9&auto=format&fit=crop&w=1400&q=80",
				Tags: []string{
					"#windy",
					"#tdd",
				},
			},
		}

		if !reflect.DeepEqual(gotPosts, wantPosts) {
			t.Errorf("got posts %+v want %+v", gotPosts, wantPosts)
		}

		NewPostRepository(connection).DeleteAll()

		gotPosts = postRepository.FindAll()
		wantPosts = []Post{}

		if !reflect.DeepEqual(gotPosts, wantPosts) {
			t.Errorf("got posts %+v want %+v", gotPosts, wantPosts)
		}
	})
}
