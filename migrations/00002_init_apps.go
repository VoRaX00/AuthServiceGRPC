package migrations

import (
	"context"
	"database/sql"
	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(UpApps, DownApps)
}

func UpApps(ctx context.Context, tx *sql.Tx) error {
	query := `CREATE TABLE IF NOT EXISTS apps (
    	id SERIAL PRIMARY KEY,
    	name_app TEXT NOT NULL UNIQUE,
    	secret TEXT NOT NULL UNIQUE
	);`

	_, err := tx.ExecContext(ctx, query)
	return err
}

func DownApps(ctx context.Context, tx *sql.Tx) error {
	query := `DROP TABLE IF EXISTS apps;`
	_, err := tx.ExecContext(ctx, query)
	return err
}
