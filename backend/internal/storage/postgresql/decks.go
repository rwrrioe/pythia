package postgresql

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/rwrrioe/pythia/backend/internal/auth"
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

func (s *DeckStorage) ListBySession(ctx context.Context, sessionId int) (*models.Deck, error) {
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

	return &deck, nil
}
