package postgresql

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/rwrrioe/pythia/backend/internal/domain/entities"
	"github.com/rwrrioe/pythia/backend/internal/storage/models"
)

type DeckStorage struct {
	pool *pgxpool.Pool
}

func NewDeckStorage(pool *pgxpool.Pool) *DeckStorage {
	return &DeckStorage{pool: pool}
}

func (s *DeckStorage) ListBySession(ctx context.Context, q Querier, sessionId int64, uid int64) (*entities.Deck, error) {
	const op = "postgresql.DeckStorage.ListBySession"

	var d models.Deck
	err := q.QueryRow(ctx,
		`SELECT id, user_id, session_id
         FROM decks
         WHERE user_id=$1 AND session_id=$2`,
		uid, int(sessionId),
	).Scan(&d.Id, &d.UserId, &d.SessionId)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrDeckNotFound
		}
		return nil, fmt.Errorf("%s:%w", op, err)
	}

	return &entities.Deck{
		Id:        d.Id,
		SessionId: d.SessionId,
	}, nil
}

func (s *DeckStorage) GetOrCreate(ctx context.Context, q Querier, sessionId, uid int64) (int, error) {
	const op = "postgresql.DeckStorage.GetOrCreate"

	var id int

	err := q.QueryRow(ctx,
		`INSERT INTO decks (user_id, session_id)
         VALUES ($1, $2)
         ON CONFLICT (user_id, session_id) DO UPDATE SET session_id=EXCLUDED.session_id
        RETURNING id
         `,
		uid, sessionId,
	).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("%s:%w", op, err)
	}
	return id, nil
}

func (s *DeckStorage) AttachFlashcard(ctx context.Context, q Querier, deckId int, flId int) error {
	const op = "postgresql.DeckStorage.AttachFlashcard"

	_, err := q.Exec(ctx,
		`INSERT INTO decks_flashcards (deck_id, flashcard_id)
         VALUES ($1, $2)
         ON CONFLICT (deck_id, flashcard_id ) DO NOTHING;
         `,
		deckId, flId,
	)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			// дубликат
			return fmt.Errorf("%s:%w", op, ErrDeckFlashcardAlreadyExists)
		}
		return fmt.Errorf("%s:%w", op, err)
	}
	return nil
}
