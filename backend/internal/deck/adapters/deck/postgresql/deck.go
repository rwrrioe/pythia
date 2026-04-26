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

func (s *Storage) ListBySession(
	ctx context.Context,
	sessionId uuid.UUID,
	uid uuid.UUID,
) (*domain.Deck, error) {
	const op = "deck.postgresql.ListBySession"

	res, err := s.q.GetDeckBySession(ctx, decksqlc.GetDeckBySessionParams{
		SessionID: pgtype.UUID{
			Bytes: sessionId,
		},
		UserID: pgtype.UUID{
			Bytes: uid,
		},
	})

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("%s:%w", op, domainerrors.ErrDeckNotFound)
		}

		return nil, fmt.Errorf("%s:%w", op, err)
	}

	return &domain.Deck{
		Id:        res.UserID.Bytes,
		SessionId: res.SessionID.Bytes,
	}, nil
}
