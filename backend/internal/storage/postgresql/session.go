package postgresql

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/rwrrioe/pythia/backend/internal/auth"
	"github.com/rwrrioe/pythia/backend/internal/storage/models"
)

type SessionStorage struct {
	conn *pgx.Conn
}

func NewSessionStorage(conn *pgx.Conn) (*SessionStorage, error) {
	const op = "storage.postgres.New"

	return &SessionStorage{conn: conn}, nil
}

func (s *SessionStorage) ListSessions(ctx context.Context) ([]models.Session, error) {
	const op = "postgresql.SessionStorage.ListSessions"

	var sessions []models.Session
	userId, ok := auth.UIDFromContext(ctx)
	if !ok {
		return nil, ErrUserNotFound
	}

	if err := s.conn.QueryRow(ctx,
		"SELECT * FROM sessions WHERE user_id=$1", userId).Scan(&sessions); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, fmt.Errorf("%s:%w", op, err)
	}

	return sessions, nil
}
