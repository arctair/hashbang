package v1

import (
	"context"
	"reflect"
	"testing"

	"github.com/arctair/go-assertutil"
	"github.com/cockroachdb/cockroach-go/v2/testserver"
	"github.com/jackc/pgx/v4/pgxpool"
)

func TestNamedTagListRepository(t *testing.T) {
	testServer, err := testserver.NewTestServer()
	defer testServer.Stop()
	assertutil.NotError(t, err)

	pool, err := pgxpool.Connect(context.Background(), testServer.PGURL().String())
	assertutil.NotError(t, err)
	assertutil.NotError(t, Migrate(pool))

	deadbeef := "deadbeef-dead-beef-dead-beefdeadbeef"

	t.Run("get empty named tag lists", func(t *testing.T) {
		got, _ := NewNamedTagListRepository(pool).FindAll([]string{"bucket"})
		want := []NamedTagList{}

		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %+v want %+v", got, want)
		}
	})

	t.Run("create named tag list", func(t *testing.T) {
		if err := NewNamedTagListRepository(pool).Create(
			"bucket",
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

		if err := NewNamedTagListRepository(pool).Create(
			"bucket",
			NamedTagList{
				ID:   "a5a5acbf-1541-4fd8-bf9a-343b75b8550f",
				Name: "tag list name",
				Tags: []string{
					"#windy",
					"#tdd",
				},
			},
		); err != nil {
			t.Fatal(err)
		}

		var got []NamedTagList
		if got, err = NewNamedTagListRepository(pool).FindAll([]string{"bucket"}); err != nil {
			t.Fatal(err)
		}
		want := []NamedTagList{
			{
				ID:   "7fe6ca35-d868-48a9-94d4-6e7f7db450ea",
				Name: "tag list name",
				Tags: []string{
					"#windy",
					"#tdd",
				},
			},
			{
				ID:   "a5a5acbf-1541-4fd8-bf9a-343b75b8550f",
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
	})

	t.Run("delete named tag list by id", func(t *testing.T) {
		if err := NewNamedTagListRepository(pool).DeleteByIds([]string{"7fe6ca35-d868-48a9-94d4-6e7f7db450ea"}); err != nil {
			t.Fatal(err)
		}

		got, _ := NewNamedTagListRepository(pool).FindAll([]string{"bucket"})
		want := []NamedTagList{
			{
				ID:   "a5a5acbf-1541-4fd8-bf9a-343b75b8550f",
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
	})

	t.Run("delete all named tag lists", func(t *testing.T) {
		if err := NewNamedTagListRepository(pool).DeleteAll(); err != nil {
			t.Fatal(err)
		}

		got, _ := NewNamedTagListRepository(pool).FindAll([]string{"bucket"})
		want := []NamedTagList{}

		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %+v want %+v", got, want)
		}
	})

	t.Run("replace named tag list by id", func(t *testing.T) {
		if err := NewNamedTagListRepository(pool).Create(
			"bucket",
			NamedTagList{
				ID:   "beefdead-d868-48a9-94d4-6e7f7db450ea",
				Name: "tag list name",
				Tags: []string{
					"#windy",
					"#tdd",
				},
			},
		); err != nil {
			t.Fatal(err)
		}

		if err := NewNamedTagListRepository(pool).Create(
			"bucket",
			NamedTagList{
				ID:   "deadbeef-d868-48a9-94d4-6e7f7db450ea",
				Name: "tag list name",
				Tags: []string{
					"#windy",
					"#tdd",
				},
			},
		); err != nil {
			t.Fatal(err)
		}

		if err := NewNamedTagListRepository(pool).ReplaceByIds(
			[]string{"beefdead-d868-48a9-94d4-6e7f7db450ea"},
			NamedTagList{
				ID:   "do not update",
				Name: "replaced",
				Tags: []string{
					"#replaced",
				},
			},
		); err != nil {
			t.Fatal(err)
		}

		got, _ := NewNamedTagListRepository(pool).FindAll([]string{"bucket"})
		want := []NamedTagList{
			{
				ID:   "beefdead-d868-48a9-94d4-6e7f7db450ea",
				Name: "replaced",
				Tags: []string{
					"#replaced",
				},
			},
			{
				ID:   "deadbeef-d868-48a9-94d4-6e7f7db450ea",
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

		if err := NewNamedTagListRepository(pool).DeleteAll(); err != nil {
			t.Fatal(err)
		}
	})

	t.Run("get named tag lists does not include results from other buckets", func(t *testing.T) {
		if err := NewNamedTagListRepository(pool).Create(
			"red",
			NamedTagList{Name: "red", ID: deadbeef},
		); err != nil {
			t.Fatal(err)
		}

		if err := NewNamedTagListRepository(pool).Create(
			"blue",
			NamedTagList{Name: "blue", ID: "9c5e7bad-b2f7-4d8b-9df9-fc0e51862d8e"},
		); err != nil {
			t.Fatal(err)
		}

		var got []NamedTagList
		if got, err = NewNamedTagListRepository(pool).FindAll([]string{"red"}); err != nil {
			t.Fatal(err)
		}
		want := []NamedTagList{
			{Name: "red", ID: deadbeef},
		}

		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %+v want %+v", got, want)
		}
	})
}
