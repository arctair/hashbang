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
		got := NewNamedTagListRepository(connection).FindAll()
		want := []NamedTagList{}

		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %+v want %+v", got, want)
		}

		NewNamedTagListRepository(connection).Create(
			NamedTagList{
				Name: "tag list name",
				Tags: []string{
					"#windy",
					"#tdd",
				},
			},
		)

		got = NewNamedTagListRepository(connection).FindAll()
		want = []NamedTagList{
			{
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

		NewNamedTagListRepository(connection).DeleteAll()

		got = NewNamedTagListRepository(connection).FindAll()
		want = []NamedTagList{}

		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %+v want %+v", got, want)
		}
	})
}
