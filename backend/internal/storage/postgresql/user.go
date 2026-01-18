package postgresql

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/rwrrioe/pythia/backend/internal/auth"
	"github.com/rwrrioe/pythia/backend/internal/storage/models"
)

type UserStorage struct {
	conn *pgx.Conn
}

func NewUserStorage(conn *pgx.Conn) (*UserStorage, error) {
	const op = "storage.postgres.New"

	return &UserStorage{conn: conn}, nil
}

func (s *UserStorage) GetUser(ctx context.Context) (*models.User, error) {
	const op = "postgresql.UserStorage.GetUser"

	var user models.User
	userId, ok := auth.UIDFromContext(ctx)
	if !ok {
		return nil, ErrUserNotFound
	}

	if err := s.conn.QueryRow(ctx,
		"SELECT * FROM users WHERE user_id=$1", userId).Scan(&user); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, fmt.Errorf("%s:%w", op, err)
	}

	return &user, nil
}
