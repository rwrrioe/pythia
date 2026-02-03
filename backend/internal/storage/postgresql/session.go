package postgresql

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rwrrioe/pythia/backend/internal/domain/entities"
	"github.com/rwrrioe/pythia/backend/internal/storage/models"
)

type SessionStorage struct {
	pool *pgxpool.Pool
}

func NewSessionStorage(pool *pgxpool.Pool) *SessionStorage {

	return &SessionStorage{pool: pool}
}

const sessionCols = `
    id, name, user_id, status, lang_id, started_at, ended_at, accuracy
`

func scanSession(row pgx.Row, m *models.Session) error {
	return row.Scan(
		&m.Id,
		&m.Name,
		&m.UserId,
		&m.Status,
		&m.Lang,
		&m.StartedAt,
		&m.EndedAt,
		&m.Accuracy,
	)
}

func (s *SessionStorage) ListSessions(ctx context.Context, q Querier, uid int64) ([]entities.Session, error) {
	const op = "postgresql.SessionStorage.ListSessions"

	rows, err := q.Query(ctx,
		`SELECT `+sessionCols+`
			FROM sessions 
			WHERE user_id=$1
			ORDER BY started_at DESC			
			`, uid)
	if err != nil {
		return nil, fmt.Errorf("%s:%w", op, err)
	}

	defer rows.Close()

	out := make([]entities.Session, 0, 16)
	for rows.Next() {
		var m models.Session

		if err := scanSession(rows, &m); err != nil {
			return nil, fmt.Errorf("%s:%w", op, err)
		}

		out = append(out, entities.Session{
			Id:        m.Id,
			UserId:    m.UserId,
			Name:      m.Name,
			Status:    m.Status,
			Language:  m.Lang,
			StartedAt: m.StartedAt,
			EndedAt:   m.EndedAt,
			Accuracy:  m.Accuracy,
		})
	}

	if rows.Err() != nil {
		return nil, fmt.Errorf("%s:%w", op, rows.Err())
	}

	return out, nil
}

func (s *SessionStorage) ListLatest(ctx context.Context, q Querier, uid int64) ([]entities.Session, error) {
	const op = "postgresql.SessionStorage.ListLatestSessions"

	rows, err := q.Query(ctx,
		`							SELECT `+sessionCols+`
			 						FROM sessions 
									WHERE user_id=$1 AND 
									ended_at >= NOW() - INTERVAL '7 days' 
			 						ORDER BY ended_at 
    								DESC LIMIT 4
			`, uid)
	if err != nil {
		return nil, fmt.Errorf("%s:%w", op, err)
	}

	defer rows.Close()

	out := make([]entities.Session, 0, 4)
	for rows.Next() {
		var m models.Session

		if err := scanSession(rows, &m); err != nil {
			return nil, fmt.Errorf("%s:%w", op, err)
		}

		out = append(out, entities.Session{
			Id:        m.Id,
			UserId:    m.UserId,
			Name:      m.Name,
			Status:    m.Status,
			Language:  m.Lang,
			StartedAt: m.StartedAt,
			EndedAt:   m.EndedAt,
			Accuracy:  m.Accuracy,
		})
	}

	if rows.Err() != nil {
		return nil, fmt.Errorf("%s:%w", op, err)
	}

	return out, nil
}

func (s *SessionStorage) GetSession(ctx context.Context, q Querier, sessionId int64, uid int64) (*entities.Session, error) {
	const op = "postgresql.SessionStorage.GetSession"

	var m models.Session
	err := scanSession(
		q.QueryRow(ctx,
			`SELECT `+sessionCols+`
             FROM sessions
             WHERE id=$1 AND user_id=$2`, sessionId, uid),
		&m)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrSessionNotFound
		}
		return nil, fmt.Errorf("%s:%w", op, err)
	}

	ss := &entities.Session{
		Id:        m.Id,
		UserId:    m.UserId,
		Name:      m.Name,
		Status:    m.Status,
		Language:  m.Lang,
		StartedAt: m.StartedAt,
		EndedAt:   m.EndedAt,
		Accuracy:  m.Accuracy,
	}
	return ss, nil
}

func (s *SessionStorage) SaveSession(ctx context.Context, q Querier, ss entities.Session, uid int64) (int, error) {
	const op = "storage.SessionStorage.SaveSession"

	var id int
	sql := `
		INSERT INTO sessions (name, user_id, status, lang_id, started_at, ended_at, accuracy)
		VALUES ($1, $2, $3, $4, $5,$6,$7)
		RETURNING id
	`

	err := q.QueryRow(ctx, sql,
		ss.Name,
		uid,
		ss.Status,
		ss.Language,
		ss.StartedAt,
		ss.EndedAt,
		ss.Accuracy).Scan(&id)

	if err != nil {
		var pgErr *pgconn.PgError

		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return 0, fmt.Errorf("%s:%w", op, ErrSessionAlreadyExists)
		}

		return 0, fmt.Errorf("%s:%w", op, err)
	}

	return id, nil
}

func (s *SessionStorage) TryMarkFinished(ctx context.Context, q Querier, sessionId int64, uid int64, endedAt time.Time) (bool, error) {
	const op = "postgresql.SessionStorage.MarkFinished"

	cmd, err := q.Exec(ctx, `
		UPDATE sessions
		SET status = 'finished',
		    ended_at =$1
		WHERE id = $2 
			AND user_id=$3
			AND status <> 'finished'
		`, endedAt, sessionId, uid)
	if err != nil {
		return true, fmt.Errorf("%s:%w", op, err)
	}

	if cmd.RowsAffected() == 0 {
		var status string

		err := q.QueryRow(ctx,
			`
				SELECT status
				FROM sessions
				WHERE id = $1 
					AND user_id= $2
			`).Scan(&status)

		if errors.Is(err, pgx.ErrNoRows) {
			return true, fmt.Errorf("%s:%w", op, ErrSessionNotFound)
		}

		if status == "finished" {
			return false, nil
		}

		return true, fmt.Errorf("%s:%w", op, err)
	}

	return true, nil
}

func (s *SessionStorage) UpdateAccuracy(ctx context.Context, q Querier, sessionId int64, uid int64, accuracy float64) error {
	const op = "postgresql.SessionStorage.UpdateAccuracy"

	cmd, err := q.Exec(ctx, `
        UPDATE sessions
        SET accuracy = $1
        WHERE id = $2 AND user_id = $3
    `, accuracy, sessionId, uid)
	if err != nil {
		return fmt.Errorf("%s:%w", op, err)
	}
	if cmd.RowsAffected() == 0 {
		return ErrSessionNotFound
	}
	return nil
}
