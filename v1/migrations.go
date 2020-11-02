package v1

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v4"
)

// Migration ...
type Migration struct {
	Index int
	Sql   string
}

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

	migrations := []Migration{
		{Index: 5, Sql: "create table named_tag_lists (\"name\" text, \"tags\" text[])"},
		{Index: 6, Sql: "alter table named_tag_lists add column id uuid primary key"},
	}

	for _, migration := range migrations {
		if migration.Index < schemaVersion {
			fmt.Printf("Skipping migration %d: %s\n", migration.Index, migration.Sql)
		} else {
			fmt.Printf("Running migration %d: %s\n", migration.Index, migration.Sql)
			_, err = connection.Exec(context.Background(), migration.Sql)
			if err != nil {
				return fmt.Errorf("Failed to migrate %d: %s", migration.Index, err)
			}
		}
	}

	_, err = connection.Exec(
		context.Background(),
		"update metadata set \"value\" = $1 where \"name\" = 'schemaVersion'",
		migrations[len(migrations)-1].Index,
	)
	if err != nil {
		return err
	}

	return err
}
