package postgresql

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/rwrrioe/pythia/backend/internal/auth"
	"github.com/rwrrioe/pythia/backend/internal/domain/entities"
	"github.com/rwrrioe/pythia/backend/internal/storage/models"
)

type SessionStorage struct {
	conn *pgx.Conn
}

func NewSessionStorage(conn *pgx.Conn) (*SessionStorage, error) {
	const op = "storage.postgres.New"

	return &SessionStorage{conn: conn}, nil
}

func (s *SessionStorage) ListSessions(ctx context.Context) ([]entities.Session, error) {
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

	var ssions []entities.Session
	for _, s := range sessions {
		ssions = append(ssions, entities.Session{
			Id:        s.Id,
			Name:      s.Name,
			Accuracy:  s.Accuracy,
			EndedAt:   s.EndedAt,
			StartedAt: s.StartedAt,
			Status:    s.Status,
			Language:  s.Lang,
		})
	}

	return ssions, nil
}

func (s *SessionStorage) ListLatest(ctx context.Context) ([]entities.Session, error) {
	const op = "postgresql.SessionStorage.ListLatestSessions"

	var sessions []models.Session
	userId, ok := auth.UIDFromContext(ctx)
	if !ok {
		return nil, ErrUserNotFound
	}

	rows, err := s.conn.Query(ctx,
		`SELECT *
			 						FROM sessions 
									WHERE user_id=$1 AND 
									ended_at >= NOW() - INTERVAL '7 days' 
			 						ORDER BY ended_at 
    								DESC LIMIT 4
			`, userId)
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

	var ssions []entities.Session
	for _, s := range sessions {
		ssions = append(ssions, entities.Session{
			Id:        s.Id,
			Name:      s.Name,
			Accuracy:  s.Accuracy,
			EndedAt:   s.EndedAt,
			StartedAt: s.StartedAt,
			Status:    s.Status,
			Language:  s.Lang,
		})
	}

	return ssions, nil
}

func (s *SessionStorage) GetSession(ctx context.Context, sessionId int) (*entities.Session, error) {
	const op = "postgresql.SessionStorage.GetSession"

	var session models.Session
	userId, ok := auth.UIDFromContext(ctx)
	if !ok {
		return nil, ErrUserNotFound
	}

	if err := s.conn.QueryRow(ctx, `SELECT *
										FROM sessions 
										WHERE session_id=$1 
										AND user_id=$2`, sessionId, userId).Scan(&session); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, fmt.Errorf("%s:%w", op, err)
	}

	return &entities.Session{
		Id:        sessionId,
		Name:      session.Name,
		Accuracy:  session.Accuracy,
		EndedAt:   session.EndedAt,
		StartedAt: session.StartedAt,
		Status:    session.Status,
		Language:  session.Lang,
	}, nil
}

func (s *SessionStorage) SaveSession(ctx context.Context, ss entities.Session) error {
	const op = "storage.SessionStorage.SaveSession"

	uid, ok := auth.UIDFromContext(ctx)
	if !ok {
		return ErrUserNotFound
	}

	sql := `
		INSERT INTO sessions (id, name, user_id, status, lang_id, started_at, ended_at, accuracy)
		VALUES ($1, $2, $3, $4, $5, $6,$7,$8)
	`

	_, err := s.conn.Exec(ctx, sql, ss.Id, ss.Name, uid, ss.Status, ss.Language, ss.StartedAt, ss.EndedAt, ss.Accuracy)
	if err != nil {
		var pgErr *pgconn.PgError

		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return fmt.Errorf("%s:%w", op, ErrSessionAlreadyExists)
		}

		return fmt.Errorf("%s:%w", op, err)
	}

	return nil
}
