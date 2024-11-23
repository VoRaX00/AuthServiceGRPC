package migrations

import (
	"context"
	"database/sql"
	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(InsertApp, DeleteApp)
}

func InsertApp(ctx context.Context, tx *sql.Tx) error {
	query := `INSERT INTO apps (id, name_app, secret) VALUES (1, 'test', 'test-secret') ON CONFLICT DO NOTHING;`
	_, err := tx.ExecContext(ctx, query)
	return err
}

func DeleteApp(ctx context.Context, tx *sql.Tx) error {
	query := `DELETE FROM apps WHERE id = 1;`
	_, err := tx.ExecContext(ctx, query)
	return err
}
