package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/lib/pq"
	"sso/internal/domain/models"
	"sso/internal/storage"
)

type Storage struct {
	db *sql.DB
}

func New(db *sql.DB) storage.Provider {
	return &Storage{db: db}
}

func (s *Storage) SaveUser(ctx context.Context, email string, passHash []byte) (int64, error) {
	const op = "storage.postgres.SaveUser"

	var id int64
	err := s.db.QueryRowContext(ctx,
		"INSERT INTO users(email, pass_hash) VALUES ($1, $2) RETURNING id",
		email, passHash,
	).Scan(&id)

	if err != nil {
		if isDuplicate(err) {
			return 0, fmt.Errorf("%s: %w", op, storage.ErrUserExists)
		}
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}

func (s *Storage) User(ctx context.Context, email string) (models.User, error) {
	const op = "storage.postgres.User"

	var user models.User
	err := s.db.QueryRowContext(ctx,
		"SELECT id, email, pass_hash FROM users WHERE email = $1",
		email,
	).Scan(&user.ID, &user.Email, &user.PassHash)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.User{}, fmt.Errorf("%s: %w", op, storage.ErrUserNotFound)
		}
		return models.User{}, fmt.Errorf("%s: %w", op, err)
	}

	return user, nil
}

func (s *Storage) IsAdmin(ctx context.Context, userID int64) (bool, error) {
	const op = "storage.postgres.IsAdmin"

	var isAdmin bool
	err := s.db.QueryRowContext(ctx,
		"SELECT is_admin FROM users WHERE id = $1",
		userID,
	).Scan(&isAdmin)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, fmt.Errorf("%s: %w", op, storage.ErrUserNotFound)
		}
		return false, fmt.Errorf("%s: %w", op, err)
	}

	return isAdmin, nil
}

func (s *Storage) App(ctx context.Context, appID int) (models.App, error) {
	const op = "storage.postgres.App"

	var app models.App
	err := s.db.QueryRowContext(ctx,
		"SELECT id, name, secret FROM apps WHERE id = $1",
		appID,
	).Scan(&app.ID, &app.Name, &app.Secret)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.App{}, fmt.Errorf("%s: %w", op, storage.ErrAppNotFound)
		}
		return models.App{}, fmt.Errorf("%s: %w", op, err)
	}

	return app, nil
}

func (s *Storage) Ping(ctx context.Context) error {
	return s.db.PingContext(ctx)
}

func isDuplicate(err error) bool {
	var pgErr *pq.Error
	return errors.As(err, &pgErr) && pgErr.Code == "23505"
}
