package v1

import (
	"context"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

// NamedTagListRepository ...
type NamedTagListRepository interface {
	FindAllOld() ([]NamedTagList, error)
	FindAll(buckets []string) ([]NamedTagList, error)
	Create(bucket string, namedTagList NamedTagList) error
	ReplaceByIds(ids []string, ntl NamedTagList) error
	DeleteAll() error
	DeleteByIds(ids []string) error
}

type namedTagListRepository struct {
	pool *pgxpool.Pool
}

func (r *namedTagListRepository) FindAllOld() ([]NamedTagList, error) {
	var (
		rows pgx.Rows
		err  error
	)

	if rows, err = r.pool.Query(context.Background(), "select \"id\", \"name\", \"tags\" from named_tag_lists"); err != nil {
		return nil, err
	}

	namedTagLists := []NamedTagList{}

	var namedTagList NamedTagList
	for rows.Next() {
		if err = rows.Scan(&namedTagList.ID, &namedTagList.Name, &namedTagList.Tags); err != nil {
			return nil, err
		}
		namedTagLists = append(namedTagLists, namedTagList)
	}

	return namedTagLists, nil
}

func (r *namedTagListRepository) FindAll(buckets []string) ([]NamedTagList, error) {
	var (
		rows pgx.Rows
		err  error
	)

	if rows, err = r.pool.Query(context.Background(), "select \"id\", \"name\", \"tags\" from named_tag_lists"); err != nil {
		return nil, err
	}

	namedTagLists := []NamedTagList{}

	var namedTagList NamedTagList
	for rows.Next() {
		if err = rows.Scan(&namedTagList.ID, &namedTagList.Name, &namedTagList.Tags); err != nil {
			return nil, err
		}
		namedTagLists = append(namedTagLists, namedTagList)
	}

	return namedTagLists, nil
}

func (r *namedTagListRepository) Create(bucket string, namedTagList NamedTagList) error {
	_, err := r.pool.Exec(
		context.Background(),
		"insert into named_tag_lists (\"id\", \"name\", \"tags\") values ($1, $2, $3)",
		namedTagList.ID,
		namedTagList.Name,
		namedTagList.Tags,
	)
	return err
}

func (r *namedTagListRepository) ReplaceByIds(ids []string, ntl NamedTagList) error {
	_, err := r.pool.Exec(
		context.Background(),
		"update named_tag_lists set \"name\" = $1, \"tags\" = $2 where \"id\" = ANY($3)",
		ntl.Name,
		ntl.Tags,
		ids,
	)
	return err
}

func (r *namedTagListRepository) DeleteAll() error {
	_, err := r.pool.Exec(context.Background(), "delete from named_tag_lists")
	return err
}

func (r *namedTagListRepository) DeleteByIds(ids []string) error {
	_, err := r.pool.Exec(
		context.Background(),
		"delete from named_tag_lists where \"id\" = ANY($1)",
		ids,
	)
	return err
}

// NewNamedTagListRepository ...
func NewNamedTagListRepository(pool *pgxpool.Pool) NamedTagListRepository {
	return &namedTagListRepository{
		pool: pool,
	}
}
