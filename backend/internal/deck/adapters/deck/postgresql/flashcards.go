package postgresql

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	decksqlc "github.com/rwrrioe/pythia/backend/internal/deck/adapters/postgresql/sqlc"
	"github.com/rwrrioe/pythia/backend/internal/deck/domain"
	"github.com/rwrrioe/pythia/backend/internal/deck/domain/errors"
)

func (s *Storage) ListByDeck(
	ctx context.Context,
	deckId uuid.UUID,
	uid uuid.UUID,
) ([]domain.FlashCard, error) {
	const op = "deck.postgresql.ListByDeck"

	res, err := s.q.ListFlashcardsByDeck(ctx, decksqlc.ListFlashcardsByDeckParams{
		UserID: pgtype.UUID{
			Bytes: uid,
		},
		DeckID: pgtype.UUID{
			Bytes: deckId,
		},
	})

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("%s:%w", op, domainerrors.ErrDeckNotFound)
		}

		return nil, fmt.Errorf("%s:%w", op, err)
	}

	var flashcards []domain.FlashCard

	for _, row := range res {
		flashcards = append(flashcards, domain.FlashCard{
			Id:     row.ID.Bytes,
			Word:   row.Word.String,
			Transl: row.Transl.String,
			Lang:   int(row.LangID.Int32),
		})
	}

	return flashcards, nil
}
