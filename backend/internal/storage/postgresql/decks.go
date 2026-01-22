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

type DeckStorage struct {
	conn *pgx.Conn
}

func NewDeckStorage(conn *pgx.Conn) *DeckStorage {
	return &DeckStorage{
		conn: conn,
	}
}

func (s *DeckStorage) ListBySession(ctx context.Context, sessionId int) (*entities.Deck, error) {
	const op = "Storage.DeckStorage.ListBySession"

	var deck models.Deck

	uid, ok := auth.UIDFromContext(ctx)
	if !ok {
		return nil, ErrUserNotFound
	}

	if err := s.conn.QueryRow(ctx, "SELECT * FROM decks WHERE user_id=$1 AND session_id=$2", uid, sessionId).Scan(&deck); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrDeckNotFound
		}
		return nil, fmt.Errorf("%s:%w", op, err)
	}

	return &entities.Deck{
		Id:        deck.Id,
		SessionId: deck.SessionId,
	}, nil
}

func (s *DeckStorage) SaveDeck(ctx context.Context, deck entities.Deck) error {
	const op = "storage.DeckStorage.SaveDeck"

	uid, ok := auth.UIDFromContext(ctx)
	if !ok {
		return ErrUserNotFound
	}

	sql := `
		INSERT INTO decks (id, user_id, session_id)
		VALUES ($1, $2, $3)
	`
	_, err := s.conn.Exec(ctx, sql, deck.Id, uid, deck.SessionId)
	if err != nil {
		var pgErr *pgconn.PgError

		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return fmt.Errorf("%s:%w", op, ErrDeckAlreadyExists)
		}

		return fmt.Errorf("%s:%w", op, err)
	}
	return nil
}

func (s *DeckStorage) AttachFlashcard(ctx context.Context, deckId int, flId int) error {
	const op = "storage.DeckStorage.AttachFlashcard"

	sql := `
		INSERT INTO decks_flashcards (deck_id, flashcard_id)
		VALUES ($1, $2)
	`
	_, err := s.conn.Exec(ctx, sql, deckId, flId)
	if err != nil {
		var pgErr *pgconn.PgError

		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return fmt.Errorf("%s:%w", op, "deck-flashcard relation already exists")
		}

		return fmt.Errorf("%s:%w", op, err)
	}
	return nil
}
