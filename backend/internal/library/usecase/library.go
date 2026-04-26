package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/rwrrioe/pythia/backend/internal/auth/authn"
	"github.com/rwrrioe/pythia/backend/internal/domain/entities"
	"github.com/rwrrioe/pythia/backend/internal/storage/postgresql"
)

type LibraryService struct {
	session SessionProvider

	pool postgresql.Querier
	txm  *postgresql.TxManager
}

func NewLibraryService(
	session SessionProvider,
	pool postgresql.Querier,
	txm *postgresql.TxManager,
) *LibraryService {
	return &LibraryService{
		session: session,
		pool:    pool,
		txm:     txm,
	}
}

func (s *LibraryService) Library(ctx context.Context) ([]entities.Session, error) {
	const op = "service.Libraryservice.Library"

	uid, ok := authn.UIDFromContext(ctx)
	if !ok {
		return nil, ErrUnauthorized
	}

	sessions, err := s.session.ListSessions(ctx, s.pool, uid)
	if err != nil {
		return nil, fmt.Errorf("%s:%w", op, err)
	}

	return sessions, nil
}

func (s *LibraryService) GetSession(ctx context.Context, sessionId uuid.UUID) (*entities.Session, error) {
	const op = "service.Libraryservice.GetSession"

	uid, ok := authn.UIDFromContext(ctx)
	if !ok {
		return nil, ErrUnauthorized
	}

	ss, err := s.session.GetSession(ctx, s.pool, sessionId, uid)
	if err != nil {
		if errors.Is(err, postgresql.ErrSessionNotFound) {
			return nil, fmt.Errorf("%s:%w", op, ErrSessionNotFound)
		}

		return nil, fmt.Errorf("%s:%w", op, err)
	}

	return ss, nil
}
