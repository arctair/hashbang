package v1

import (
	"errors"
	"fmt"
	"reflect"
	"testing"
)

type stubNamedTagListRepositoryForService struct {
	NamedTagListRepository

	withBucket       string
	withNamedTagList NamedTagList
	willError        bool

	err error
}

func (r *stubNamedTagListRepositoryForService) CreateOld(namedTagList NamedTagList) error {
	return errors.New("do not call")
}

func (r *stubNamedTagListRepositoryForService) Create(bucket string, namedTagList NamedTagList) error {
	requestMatched := bucket == r.withBucket && reflect.DeepEqual(namedTagList, r.withNamedTagList)
	if !requestMatched {
		r.err = fmt.Errorf("Stub got bucket %s want %s got named tag list %+v want %+v", bucket, r.withBucket, namedTagList, r.withNamedTagList)
	}
	if requestMatched == r.willError {
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
		repository := &stubNamedTagListRepositoryForService{
			withBucket:       "bucket",
			withNamedTagList: response,
		}
		service := NewNamedTagListService(
			repository,
			&stubUUIDGenerator{
				response: "3e99aa77-615e-4a55-930d-d4c77cfd1b72",
			},
		)

		var (
			gotResponse *NamedTagList
			err         error
		)
		if gotResponse, err = service.Create("bucket", request); err != nil {
			t.Fatal(err)
		}

		if repository.err != nil {
			t.Error(repository.err)
		}

		if !reflect.DeepEqual(gotResponse, &response) {
			t.Errorf("got %+v want %+v", gotResponse, &response)
		}
	})

	t.Run("create when repository has error", func(t *testing.T) {
		repository := &stubNamedTagListRepositoryForService{
			withBucket:       "bucket",
			withNamedTagList: request,
			willError:        true,
		}
		service := NewNamedTagListService(
			repository,
			&stubUUIDGenerator{},
		)

		_, gotErr := service.Create("bucket", request)

		if repository.err != nil {
			t.Error(repository.err)
		}

		if gotErr == nil {
			t.Fatal("got no error")
		}

		wantErr := "there was an error"

		if gotErr.Error() != wantErr {
			t.Errorf("got error %s want %s", gotErr.Error(), wantErr)
		}
	})
}
