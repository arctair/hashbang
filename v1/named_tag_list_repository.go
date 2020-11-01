package v1

import (
	"context"

	"github.com/jackc/pgx/v4"
)

// NamedTagListRepository ...
type NamedTagListRepository interface {
	FindAll() []NamedTagList
	Create(namedTagList NamedTagList)
	DeleteAll()
}

type namedTagListRepository struct {
	connection *pgx.Conn
}

func (r *namedTagListRepository) FindAll() []NamedTagList {
	rows, err := r.connection.Query(context.Background(), "select \"name\", \"tags\" from named_tag_lists")
	if err != nil {
		panic(err)
	}

	namedTagLists := []NamedTagList{}

	var namedTagList NamedTagList
	for rows.Next() {
		rows.Scan(&namedTagList.Name, &namedTagList.Tags)
		namedTagLists = append(namedTagLists, namedTagList)
	}

	return namedTagLists
}

func (r *namedTagListRepository) Create(namedTagList NamedTagList) {
	_, err := r.connection.Exec(context.Background(), "insert into named_tag_lists (\"name\", \"tags\") values ($1, $2)", namedTagList.Name, namedTagList.Tags)
	if err != nil {
		panic(err)
	}
}

func (r *namedTagListRepository) DeleteAll() {
	_, err := r.connection.Exec(context.Background(), "delete from named_tag_lists")
	if err != nil {
		panic(err)
	}
}

// NewNamedTagListRepository ...
func NewNamedTagListRepository(connection *pgx.Conn) NamedTagListRepository {
	return &namedTagListRepository{
		connection: connection,
	}
}
