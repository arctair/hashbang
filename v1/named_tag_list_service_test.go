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
	if r.willError || !reflect.DeepEqual(namedTagList, r.request) {
		return errors.New("there was an error")
	}
	return nil
}

type stubUUIDGenerator struct {
	response string
}

func (r *stubUUIDGenerator) Generate() string {
	return r.response
}

func TestNamedTagListService(t *testing.T) {
	request := NamedTagList{
		Name: "tag list name",
		Tags: []string{
			"#windy",
			"#tdd",
		},
	}
	t.Run("create", func(t *testing.T) {
		response := NamedTagList{
			ID:   "3e99aa77-615e-4a55-930d-d4c77cfd1b72",
			Name: "tag list name",
			Tags: []string{
				"#windy",
				"#tdd",
			},
		}
		service := NewNamedTagListService(
			&stubNamedTagListRepositoryForService{
				request: response,
			},
			&stubUUIDGenerator{
				response: "3e99aa77-615e-4a55-930d-d4c77cfd1b72",
			},
		)

		gotResponse, err := service.CreateOld(request)
		if err != nil {
			t.Fatal(err)
		}

		if !reflect.DeepEqual(gotResponse, &response) {
			t.Errorf("got %+v want %+v", gotResponse, &response)
		}
	})

	t.Run("create when repository has error", func(t *testing.T) {
		service := NewNamedTagListService(
			&stubNamedTagListRepositoryForService{
				request:   request,
				willError: true,
			},
			&stubUUIDGenerator{},
		)

		_, gotErr := service.CreateOld(request)

		if gotErr == nil {
			t.Fatal("got no error")
		}

		wantErr := "there was an error"

		if gotErr.Error() != wantErr {
			t.Errorf("got error %s want %s", gotErr.Error(), wantErr)
		}
	})
}
