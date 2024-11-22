package migrations

import (
	"context"
	"database/sql"
	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(AddColumnIsAdmin, DropColumnIsAdmin)
}

func AddColumnIsAdmin(ctx context.Context, tx *sql.Tx) error {
	query := `ALTER TABLE users ADD COLUMN is_admin BOOLEAN NOT NULL DEFAULT FALSE;`
	_, err := tx.ExecContext(ctx, query)
	return err
}

func DropColumnIsAdmin(ctx context.Context, tx *sql.Tx) error {
	query := `ALTER TABLE users DROP COLUMN is_admin;`
	_, err := tx.ExecContext(ctx, query)
	return err
}
