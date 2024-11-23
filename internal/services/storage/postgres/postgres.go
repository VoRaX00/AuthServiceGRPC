package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"sso/internal/domain/models"
	"sso/internal/services/storage"
)

type Storage struct {
	db *sqlx.DB
}

func New(storagePath string) (*Storage, error) {
	const op = "storage.postgres.New"

	db, err := sqlx.Open("postgres", storagePath)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return &Storage{db: db}, nil
}

const ErrConstraintUnique = "23505"

// SaveUser saves the users
func (s *Storage) SaveUser(ctx context.Context, email string, passHash []byte) (int64, error) {
	const op = "storage.postgres.SaveUser"

	stmt, err := s.db.Prepare(`INSERT INTO users (email, pass_hash) VALUES (?, ?)`)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	defer func() {
		_ = stmt.Close()
	}()

	res, err := stmt.ExecContext(ctx, email, passHash)
	if err != nil {
		var postgresErr *pq.Error
		if errors.As(err, &postgresErr) && postgresErr.Code == ErrConstraintUnique {
			return 0, fmt.Errorf("%s: %w", op, storage.ErrUserExists)
		}
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	return id, nil
}

// User returns user by email
func (s *Storage) User(ctx context.Context, email string) (models.User, error) {
	const op = "storage.postgres.User"

	stmt, err := s.db.Prepare(`SELECT id, email, pass_hash FROM users WHERE email = ?`)
	if err != nil {
		return models.User{}, fmt.Errorf("%s: %w", op, err)
	}

	defer func() {
		_ = stmt.Close()
	}()

	var user models.User
	row := stmt.QueryRowContext(ctx, email)
	if err = row.Scan(&user); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.User{}, storage.ErrUserNotFound
		}
		return models.User{}, fmt.Errorf("%s: %w", op, err)
	}
	return user, nil
}

// IsAdmin checks the admin user or not
func (s *Storage) IsAdmin(ctx context.Context, userID int64) (bool, error) {
	const op = "storage.postgres.IsAdmin"

	stmt, err := s.db.Prepare(`SELECT EXISTS (SELECT 1 FROM users WHERE id = ?)`)
	if err != nil {
		return false, fmt.Errorf("%s: %w", op, err)
	}

	defer func() {
		_ = stmt.Close()
	}()

	var exists bool
	err = stmt.QueryRowContext(ctx, userID).Scan(&exists)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, fmt.Errorf("%s: %w", op, storage.ErrUserNotFound)
		}
		return false, fmt.Errorf("%s: %w", op, err)
	}
	return exists, nil
}

// App returns app by id
func (s *Storage) App(ctx context.Context, appID int32) (models.App, error) {
	const op = "storage.postgres.App"
	stmt, err := s.db.Prepare(`SELECT id, name_app, secret FROM apps WHERE id = ?`)
	if err != nil {
		return models.App{}, fmt.Errorf("%s: %w", op, err)
	}

	defer func() {
		_ = stmt.Close()
	}()

	var app models.App
	row := stmt.QueryRowContext(ctx, appID)
	if err = row.Scan(&app.ID, &app.Name, &app.Secret); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.App{}, storage.ErrAppNotFound
		}
		return models.App{}, fmt.Errorf("%s: %w", op, err)
	}
	return app, nil
}
