package v1

import (
	"context"

	"github.com/jackc/pgx/v4"
)

// Migrate ...
func Migrate(connection *pgx.Conn) error {
	_, err := connection.Exec(context.Background(), "create table posts (\"imageUri\" text, \"tags\" text[])")
	return err
}
