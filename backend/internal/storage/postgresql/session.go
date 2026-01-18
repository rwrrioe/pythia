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

	rows, err := s.conn.Query(ctx,
		"SELECT * FROM sessions WHERE user_id=$1", userId)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, fmt.Errorf("%s:%w", op, err)
	}

	defer rows.Close()

	for rows.Next() {
		var session models.Session

		if err = rows.Scan(&session); err != nil {
			return nil, fmt.Errorf("%s:%w", op, err)
		}

		sessions = append(sessions, session)
	}

	if rows.Err() != nil {
		return nil, fmt.Errorf("%s:%w", op, err)
	}

	return sessions, nil
}

func (s *SessionStorage) ListLatest(ctx context.Context) ([]models.Session, error) {
	const op = "postgresql.SessionStorage.ListLatestSessions"

	var sessions []models.Session
	userId, ok := auth.UIDFromContext(ctx)
	if !ok {
		return nil, ErrUserNotFound
	}

	rows, err := s.conn.Query(ctx,
		"SELECT * FROM sessions WHERE user_id=$1 AND ended_at >= NOW() - INTERVAL '7 days' ORDER BY ended_at DESC LIMIT 4", userId)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, fmt.Errorf("%s:%w", op, err)
	}

	defer rows.Close()

	for rows.Next() {
		var session models.Session

		if err = rows.Scan(&session); err != nil {
			return nil, fmt.Errorf("%s:%w", op, err)
		}

		sessions = append(sessions, session)
	}

	if rows.Err() != nil {
		return nil, fmt.Errorf("%s:%w", op, err)
	}

	return sessions, nil
}

func (s *SessionStorage) GetSession(ctx context.Context, sessionId int) (*models.Session, error) {
	const op = "postgresql.SessionStorage.GetSession"

	var session models.Session
	userId, ok := auth.UIDFromContext(ctx)
	if !ok {
		return nil, ErrUserNotFound
	}

	if err := s.conn.QueryRow(ctx, "SELECT * FROM sessions WHERE session_id=$1 AND user_id=$2", sessionId, userId).Scan(&session); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, fmt.Errorf("%s:%w", op, err)
	}

	return &session, nil
}
