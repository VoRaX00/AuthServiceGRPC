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
