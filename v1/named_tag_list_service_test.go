package v1

import (
	"errors"
	"reflect"
	"testing"
)

type stubNamedTagListRepositoryForService struct {
	NamedTagListRepository

	request   NamedTagList
	willError bool
}

func (r *stubNamedTagListRepositoryForService) Create(namedTagList NamedTagList) error {
	if r.willError && reflect.DeepEqual(namedTagList, r.request) {
		return errors.New("there was an error")
	}
	return nil
}

func TestNamedTagListService(t *testing.T) {
	dummyNamedTagList := NamedTagList{
		Name: "tag list name",
		Tags: []string{
			"#windy",
			"#tdd",
		},
	}

	t.Run("create", func(t *testing.T) {
		service := NewNamedTagListService(
			&stubNamedTagListRepositoryForService{
				request: dummyNamedTagList,
			},
		)

		got, err := service.Create(dummyNamedTagList)
		if err != nil {
			t.Fatal(err)
		}
		want := &dummyNamedTagList

		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %+v want %+v", got, want)
		}
	})

	t.Run("create when repository has error", func(t *testing.T) {
		service := NewNamedTagListService(
			&stubNamedTagListRepositoryForService{
				request:   dummyNamedTagList,
				willError: true,
			},
		)

		_, gotErr := service.Create(dummyNamedTagList)

		if gotErr == nil {
			t.Fatal("got no error")
		}

		wantErr := "there was an error"

		if gotErr.Error() != wantErr {
			t.Errorf("got error %s want %s", gotErr.Error(), wantErr)
		}
	})
}
