package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/rwrrioe/pythia/backend/internal/auth/authn"
	"github.com/rwrrioe/pythia/backend/internal/domain/entities"
	"github.com/rwrrioe/pythia/backend/internal/storage/postgresql"
)

type LibraryService struct {
	session    SessionProvider
	flashcards FlashCardProvider

	txm *postgresql.TxManager
}

func NewLibraryService(
	session SessionProvider,
	txm *postgresql.TxManager,
) *LibraryService {
	return &LibraryService{
		session: session,
		txm:     txm,
	}
}

func (s *LibraryService) Library(ctx context.Context) ([]entities.Session, error) {
	const op = "service.Libraryservice.Library"

	uid, ok := authn.UIDFromContext(ctx)
	if !ok {
		return nil, ErrUnauthorized
	}

	q := postgresql.NewPoolQuerier(s.txm.Pool)
	sessions, err := s.session.ListSessions(ctx, q, uid)
	if err != nil {
		return nil, fmt.Errorf("%s:%w", op, err)
	}

	return sessions, nil
}

func (s *LibraryService) GetSession(ctx context.Context, sessionId int64) (*entities.Session, error) {
	const op = "service.Libraryservice.GetSession"

	uid, ok := authn.UIDFromContext(ctx)
	if !ok {
		return nil, ErrUnauthorized
	}

	q := postgresql.NewPoolQuerier(s.txm.Pool)
	ss, err := s.session.GetSession(ctx, q, sessionId, uid)
	if err != nil {
		if errors.Is(err, postgresql.ErrSessionNotFound) {
			return nil, fmt.Errorf("%s:%w", op, ErrSessionNotFound)
		}

		return nil, fmt.Errorf("%s:%w", op, err)
	}

	return ss, nil
}
