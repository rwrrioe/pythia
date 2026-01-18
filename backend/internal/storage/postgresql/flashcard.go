package postgresql

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/rwrrioe/pythia/backend/internal/auth"
	"github.com/rwrrioe/pythia/backend/internal/storage/models"
)

type FlashCardStorage struct {
	conn *pgx.Conn
}

func NewFlashcardRepo(conn *pgx.Conn) *FlashCardStorage {
	return &FlashCardStorage{
		conn: conn,
	}
}

func (s *FlashCardStorage) ListByDeck(ctx context.Context, deckId int) ([]models.FlashCard, error) {
	const op = "Storage.FlashCardStorage.ListByDeck"

	var flashcards []models.FlashCard
	uid, ok := auth.UIDFromContext(ctx)
	if !ok {
		return nil, ErrUserNotFound
	}

	rows, err := s.conn.Query(ctx, `SELECT f.id, f.word, f.trans, f.description, l.language 
										FROM decks_flashcards df 
										JOIN flashcards f ON df.flashcard_id = f.id
										JOIN languages l ON f.lang_id = l.id
										WHERE f.user_id = $1 AND df.deck_id = $2
										`, uid, deckId)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrDeckNotFound
		}
		return nil, fmt.Errorf("%s:%w", op, err)
	}
	defer rows.Close()

	for rows.Next() {
		var fc models.FlashCard

		if err = rows.Scan(&fc); err != nil {
			return nil, fmt.Errorf("%s:%w", op, err)
		}

		flashcards = append(flashcards, fc)
	}

	if rows.Err() != nil {
		return nil, fmt.Errorf("%s:%w", op, err)
	}
	return flashcards, nil
}
