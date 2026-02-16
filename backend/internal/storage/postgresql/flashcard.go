package postgresql

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/rwrrioe/pythia/backend/internal/domain/entities"
	"github.com/rwrrioe/pythia/backend/internal/storage/models"
)

type FlashCardStorage struct {
	pool *pgxpool.Pool
}

func NewFlashcardStorage(pool *pgxpool.Pool) *FlashCardStorage {
	return &FlashCardStorage{pool: pool}
}

func scanFlashcard(row pgx.Row, m *models.FlashCard) error {
	return row.Scan(
		&m.Id,
		&m.Word,
		&m.Transl,
		&m.Lang,
	)
}

func (s *FlashCardStorage) FlashcardsPool() *pgxpool.Pool {
	return s.pool
}

// flashcards конкретной деки
func (s *FlashCardStorage) ListByDeck(ctx context.Context, q Querier, deckId uuid.UUID, uid int64) ([]entities.FlashCard, error) {
	const op = "postgresql.FlashCardStorage.ListByDeck"

	rows, err := q.Query(ctx,
		`SELECT f.id, f.word, f.transl, f.lang_id
         FROM decks_flashcards df 
         JOIN flashcards f ON df.flashcard_id = f.id
         WHERE f.user_id=$1 AND df.deck_id=$2
         ORDER BY f.id`,
		uid, deckId,
	)
	if err != nil {
		return nil, fmt.Errorf("%s:%w", op, err)
	}
	defer rows.Close()

	out := make([]entities.FlashCard, 0, 16) //todo !! маппинг ошибок добавить

	for rows.Next() {
		var m models.FlashCard
		if err := scanFlashcard(rows, &m); err != nil {
			return nil, fmt.Errorf("%s:%w", op, err)
		}

		out = append(out, entities.FlashCard{
			Id:     m.Id,
			Word:   m.Word,
			Transl: m.Transl,
			Lang:   m.Lang,
			Desc:   "", // description нет в БД
		})
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s:%w", op, err)
	}

	return out, nil
}

// flashcards пользователя
func (s *FlashCardStorage) List(ctx context.Context, q Querier, uid int64) ([]entities.FlashCard, error) {
	const op = "postgresql.FlashCardStorage.List"

	rows, err := q.Query(ctx,
		`SELECT id, word, transl, lang_id
         FROM flashcards
         WHERE user_id=$1
         ORDER BY id DESC`,
		uid,
	)
	if err != nil {
		return nil, fmt.Errorf("%s:%w", op, err)
	}
	defer rows.Close()

	out := make([]entities.FlashCard, 0, 64)

	for rows.Next() {
		var m models.FlashCard

		if err := scanFlashcard(rows, &m); err != nil {
			return nil, fmt.Errorf("%s:%w", op, err)
		}

		out = append(out, entities.FlashCard{
			Id:     m.Id,
			Word:   m.Word,
			Transl: m.Transl,
			Lang:   m.Lang,
			Desc:   "",
		})
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s:%w", op, err)
	}

	return out, nil
}

func (s *FlashCardStorage) GetOrCreate(ctx context.Context, q Querier, flCard entities.FlashCard, uid int64) (uuid.UUID, error) {
	const op = "postgresql.FlashCardStorage.GetOrCreate"

	var id uuid.UUID

	err := q.QueryRow(ctx, `
		INSERT INTO flashcards (user_id, word, transl, lang_id)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (user_id, lang_id, word)
		DO UPDATE SET transl = EXCLUDED.transl
		RETURNING id
	`, uid, flCard.Word, flCard.Transl, flCard.Lang).Scan(&id)

	if err != nil {
		return uuid.Nil, fmt.Errorf("%s:%w", op, err)
	}

	return id, nil
}
