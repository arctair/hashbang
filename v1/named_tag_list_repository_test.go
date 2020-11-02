package v1

import (
	"context"
	"reflect"
	"testing"

	"github.com/arctair/go-assertutil"
	"github.com/cockroachdb/cockroach-go/v2/testserver"
	"github.com/jackc/pgx/v4"
)

func TestNamedTagListRepository(t *testing.T) {
	testServer, err := testserver.NewTestServer()
	defer testServer.Stop()
	assertutil.NotError(t, err)

	connection, err := pgx.Connect(context.Background(), testServer.PGURL().String())
	assertutil.NotError(t, err)
	assertutil.NotError(t, Migrate(connection))

	t.Run("create, get, delete named tag list", func(t *testing.T) {
		got, _ := NewNamedTagListRepository(connection).FindAll()
		want := []NamedTagList{}

		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %+v want %+v", got, want)
		}

		if err := NewNamedTagListRepository(connection).Create(
			NamedTagList{
				ID:   "7fe6ca35-d868-48a9-94d4-6e7f7db450ea",
				Name: "tag list name",
				Tags: []string{
					"#windy",
					"#tdd",
				},
			},
		); err != nil {
			t.Fatal(err)
		}

		if got, err = NewNamedTagListRepository(connection).FindAll(); err != nil {
			t.Fatal(err)
		}
		want = []NamedTagList{
			{
				ID:   "7fe6ca35-d868-48a9-94d4-6e7f7db450ea",
				Name: "tag list name",
				Tags: []string{
					"#windy",
					"#tdd",
				},
			},
		}

		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %+v want %+v", got, want)
		}

		if err := NewNamedTagListRepository(connection).DeleteAll(); err != nil {
			t.Fatal(err)
		}

		got, _ = NewNamedTagListRepository(connection).FindAll()
		want = []NamedTagList{}

		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %+v want %+v", got, want)
		}
	})
}
