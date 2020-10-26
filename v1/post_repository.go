package v1

import (
	"context"

	"github.com/jackc/pgx/v4"
)

// PostRepository ...
type PostRepository interface {
	FindAll() []Post
	Create(post Post)
	DeleteAll()
}

type postRepository struct {
	connection *pgx.Conn
	posts      []Post
}

func (r *postRepository) FindAll() []Post {
	rows, err := r.connection.Query(context.Background(), "select \"imageUri\", \"tags\" from posts")
	if err != nil {
		panic(err)
	}

	posts := []Post{}

	var post Post
	for rows.Next() {
		rows.Scan(&post.ImageUri, &post.Tags)
		posts = append(posts, post)
	}

	return posts
}

func (r *postRepository) Create(post Post) {
	_, err := r.connection.Exec(context.Background(), "insert into posts (\"imageUri\", \"tags\") values ($1, $2)", post.ImageUri, post.Tags)
	if err != nil {
		panic(err)
	}
}

func (r *postRepository) DeleteAll() {
	_, err := r.connection.Exec(context.Background(), "delete from posts")
	if err != nil {
		panic(err)
	}
}

// NewPostRepository ...
func NewPostRepository(connection *pgx.Conn) PostRepository {
	return &postRepository{
		connection: connection,
		posts:      []Post{},
	}
}
