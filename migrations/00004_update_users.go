package migrations

import (
	"context"
	"database/sql"
	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(UpUsers2, DownUsers2)
}

func UpUsers2(ctx context.Context, tx *sql.Tx) error {
	query := `ALTER TABLE users ADD COLUMN IF NOT EXISTS pass_hash_tmp bytea;
	UPDATE users SET pass_hash_tmp = pass_hash::bytea;
	ALTER TABLE users DROP COLUMN IF EXISTS pass_hash;
	ALTER TABLE users RENAME COLUMN pass_hash_tmp TO pass_hash;`
	_, err := tx.ExecContext(ctx, query)
	return err
}

func DownUsers2(ctx context.Context, tx *sql.Tx) error {
	query := `ALTER TABLE users ADD COLUMN IF NOT EXISTS pass_hash_tmp text;
	UPDATE users SET pass_hash_tmp = pass_hash::text;
	ALTER TABLE users DROP COLUMN pass_hash;
	ALTER TABLE users RENAME COLUMN pass_hash_tmp TO pass_hash;`
	_, err := tx.ExecContext(ctx, query)
	return err
}
