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

type FlashCardStorage struct {
	conn *pgx.Conn
}

func NewFlashcardStorage(conn *pgx.Conn) *FlashCardStorage {
	return &FlashCardStorage{
		conn: conn,
	}
}

func (s *FlashCardStorage) ListByDeck(ctx context.Context, deckId int) ([]entities.FlashCard, error) {
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

	var flcards []entities.FlashCard
	for _, fl := range flashcards {
		flcards = append(flcards, entities.FlashCard{
			Id:     fl.Id,
			Word:   fl.Word,
			Transl: fl.Transl,
			Desc:   fl.Desc,
			Lang:   fl.Lang,
		})
	}

	return flcards, nil
}

func (s *FlashCardStorage) List(ctx context.Context) ([]entities.FlashCard, error) {
	const op = "Storage.FlashCardStorage.List"

	var flashcards []models.FlashCard
	uid, ok := auth.UIDFromContext(ctx)
	if !ok {
		return nil, ErrUserNotFound
	}

	rows, err := s.conn.Query(ctx, `SELECT *
										FROM flashcards
										WHERE user_id=$1
										`, uid)
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

	var flcards []entities.FlashCard
	for _, fl := range flashcards {
		flcards = append(flcards, entities.FlashCard{
			Id:     fl.Id,
			Word:   fl.Word,
			Transl: fl.Transl,
			Desc:   fl.Desc,
			Lang:   fl.Lang,
		})
	}

	return flcards, nil
}

func (s *FlashCardStorage) SaveFlashcard(ctx context.Context, flCard entities.FlashCard) error {
	const op = "storage.FlashCardStorage.SaveFlashcard"

	uid, ok := auth.UIDFromContext(ctx)
	if !ok {
		return ErrUserNotFound
	}

	sql := `
		INSERT INTO flashcards (id, user_id, word, transl ,lang_id)
		VALUES ($1, $2, $3, $4, $5)
	`

	_, err := s.conn.Exec(ctx, sql, flCard.Id, uid, flCard.Word, flCard.Transl, flCard.Lang)
	if err != nil {
		var pgErr *pgconn.PgError

		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return fmt.Errorf("%s:%w", op, ErrSessionAlreadyExists)
		}

		return fmt.Errorf("%s:%w", op, err)
	}

	return nil
}
