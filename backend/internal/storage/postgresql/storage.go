package postgresql

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	ErrUserNotFound               = errors.New("user not found")
	ErrFlashcardNotFound          = errors.New("flashcard not found")
	ErrDeckNotFound               = errors.New("deck not found")
	ErrSessionNotFound            = errors.New("session not found")
	ErrSessionAlreadyExists       = errors.New("session already exists")
	ErrSessionAlreadyFinished     = errors.New("session already finished")
	ErrNoSessions                 = errors.New("no sessions")
	ErrDeckAlreadyExists          = errors.New("deck already exists")
	ErrFlashcardAlreadyExists     = errors.New("flashcard already exists")
	ErrDeckFlashcardAlreadyExists = errors.New("deck-flashcards already exists")
)

type Querier interface {
	Exec(context.Context, string, ...any) (pgconn.CommandTag, error)
	Query(context.Context, string, ...any) (pgx.Rows, error)
	QueryRow(context.Context, string, ...any) pgx.Row
}

type PoolQuerier struct {
	pool *pgxpool.Pool
}

// / todo!!! -> pgc , костыли
func NewPoolQuerier(pool *pgxpool.Pool) *PoolQuerier {
	return &PoolQuerier{pool: pool}
}

func (q *PoolQuerier) Exec(ctx context.Context, sql string, args ...any) (pgconn.CommandTag, error) {
	return q.pool.Exec(ctx, sql, args...)
}
func (q *PoolQuerier) Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error) {
	return q.pool.Query(ctx, sql, args...)
}
func (q *PoolQuerier) QueryRow(ctx context.Context, sql string, args ...any) pgx.Row {
	return q.pool.QueryRow(ctx, sql, args...)
}

func New(ctx context.Context) (*pgxpool.Pool, error) {
	const op = "storage.postgres.New"

	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASS"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"))
	pool, err := pgxpool.New(ctx, connStr)
	if err != nil {
		return nil, fmt.Errorf("%s:%w", op, err)
	}

	return pool, nil
}
