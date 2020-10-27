package v1

import (
	"context"

	"github.com/jackc/pgx/v4"
)

// PostRepository ...
type PostRepository interface {
	FindAll() []Post
	Create(post Post) Post
	DeleteAll()
}

type postRepository struct {
	connection    *pgx.Conn
	uuidGenerator UuidGenerator
}

func (r *postRepository) FindAll() []Post {
	rows, err := r.connection.Query(context.Background(), "select \"id\", \"imageUri\", \"tags\" from posts")
	if err != nil {
		panic(err)
	}

	posts := []Post{}

	var post Post
	for rows.Next() {
		rows.Scan(&post.Id, &post.ImageUri, &post.Tags)
		posts = append(posts, post)
	}

	return posts
}

func (r *postRepository) Create(post Post) Post {
	post.Id = r.uuidGenerator.Generate()
	_, err := r.connection.Exec(context.Background(), "insert into posts (\"id\", \"imageUri\", \"tags\") values ($1, $2, $3)", post.Id, post.ImageUri, post.Tags)
	if err != nil {
		panic(err)
	}
	return post
}

func (r *postRepository) DeleteAll() {
	_, err := r.connection.Exec(context.Background(), "delete from posts")
	if err != nil {
		panic(err)
	}
}

// NewPostRepository ...
func NewPostRepository(
	connection *pgx.Conn,
	uuidGenerator UuidGenerator,
) PostRepository {
	return &postRepository{
		connection:    connection,
		uuidGenerator: uuidGenerator,
	}
}
