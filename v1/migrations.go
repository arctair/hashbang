package v1

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v4"
)

// Migrate ...
func Migrate(connection *pgx.Conn) error {
	var err error
	_, err = connection.Exec(context.Background(), "create table if not exists metadata (\"name\" text primary key, \"value\" int)")
	if err != nil {
		return err
	}

	_, err = connection.Exec(context.Background(), "insert into metadata (\"name\", \"value\") values ('schemaVersion', 0) on conflict do nothing")
	if err != nil {
		return err
	}

	row := connection.QueryRow(context.Background(), "select \"value\" from metadata where \"name\" = 'schemaVersion'")
	var schemaVersion int
	err = row.Scan(&schemaVersion)

	migrations := []string{
		"create table if not exists posts (\"imageUri\" text, \"tags\" text[])",
		"create table named_tag_lists (\"imageUri\" text, \"tags\" text[])",
		"drop table posts",
		"create table named_tag_lists_2 (\"name\" text, \"tags\" text[])",
	}

	for index, migration := range migrations {
		if index < schemaVersion {
			fmt.Printf("Skipping migration %d: %s\n", index, migration)
		} else {
			fmt.Printf("Running migration %d: %s\n", index, migration)
			_, err = connection.Exec(context.Background(), migration)
			if err != nil {
				return fmt.Errorf("Failed to migrate %d: %s", index, err)
			}
		}
	}

	_, err = connection.Exec(context.Background(), "update metadata set \"value\" = $1 where \"name\" = 'schemaVersion'", len(migrations))
	if err != nil {
		return err
	}

	return err
}
