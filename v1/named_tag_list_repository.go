package v1

import (
	"context"

	"github.com/jackc/pgx/v4"
)

// NamedTagListRepository ...
type NamedTagListRepository interface {
	FindAll() ([]NamedTagList, error)
	Create(namedTagList NamedTagList) error
	DeleteAll() error
}

type namedTagListRepository struct {
	connection *pgx.Conn
}

func (r *namedTagListRepository) FindAll() ([]NamedTagList, error) {
	var (
		rows pgx.Rows
		err  error
	)

	if rows, err = r.connection.Query(context.Background(), "select \"name\", \"tags\" from named_tag_lists"); err != nil {
		return nil, err
	}

	namedTagLists := []NamedTagList{}

	var namedTagList NamedTagList
	for rows.Next() {
		if err = rows.Scan(&namedTagList.Name, &namedTagList.Tags); err != nil {
			return nil, err
		}
		namedTagLists = append(namedTagLists, namedTagList)
	}

	return namedTagLists, nil
}

func (r *namedTagListRepository) Create(namedTagList NamedTagList) error {
	_, err := r.connection.Exec(context.Background(), "insert into named_tag_lists (\"name\", \"tags\") values ($1, $2)", namedTagList.Name, namedTagList.Tags)
	return err
}

func (r *namedTagListRepository) DeleteAll() error {
	_, err := r.connection.Exec(context.Background(), "delete from named_tag_lists")
	return err
}

// NewNamedTagListRepository ...
func NewNamedTagListRepository(connection *pgx.Conn) NamedTagListRepository {
	return &namedTagListRepository{
		connection: connection,
	}
}
