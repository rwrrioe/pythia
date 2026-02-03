package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/rwrrioe/pythia/backend/internal/auth/authn"
	"github.com/rwrrioe/pythia/backend/internal/domain/entities"
	service "github.com/rwrrioe/pythia/backend/internal/services/errors"
	"github.com/rwrrioe/pythia/backend/internal/storage/postgresql"
)

type StatsService struct {
	sessionProvider SessionProvider
	deckProvider    DeckProvider
	flCardProvider  FlashCardProvider
	txm             *postgresql.TxManager
}

func NewStatsService(
	ss SessionProvider,
	dck DeckProvider,
	fl FlashCardProvider,
	txm *postgresql.TxManager,
) *StatsService {
	return &StatsService{
		sessionProvider: ss,
		deckProvider:    dck,
		flCardProvider:  fl,
		txm:             txm,
	}
}

func (s *StatsService) Dashboard(ctx context.Context) (*entities.Dashboard, error) {
	const op = "service.StatsService.Dashboard"

	uid, ok := authn.UIDFromContext(ctx)
	if !ok {
		return nil, fmt.Errorf("%s:%w", op, service.ErrUnauthorized)
	}

	latestSessions, err := s.sessionProvider.ListLatest(ctx, s.txm.Pool, uid)
	if err != nil {
		if errors.Is(err, postgresql.ErrUserNotFound) {
			return nil, fmt.Errorf("%s:%w", op, service.ErrUnauthorized)
		}

		return nil, fmt.Errorf("%s:%w", op, err)
	}

	words, err := s.flCardProvider.List(ctx, s.txm.Pool, uid)
	if err != nil {
		return nil, fmt.Errorf("%s:%w", op, err)
	}

	sessions, err := s.sessionProvider.ListSessions(ctx, s.txm.Pool, uid)
	if err != nil {
		if errors.Is(err, postgresql.ErrUserNotFound) {
			return nil, fmt.Errorf("%s:%w", op, service.ErrUnauthorized)
		}

		return nil, fmt.Errorf("%s:%w", op, err)
	}

	if len(sessions) == 0 {
		return &entities.Dashboard{
			Streak:         0,
			WordsLearned:   0,
			Accuracy:       0,
			LatestSessions: latestSessions,
		}, nil
	}

	avgAcc := 0
	for _, s := range sessions {
		avgAcc += int(s.Accuracy)
	}
	avgAcc = avgAcc / len(sessions)

	return &entities.Dashboard{
		Streak:         0,
		WordsLearned:   len(words),
		Accuracy:       avgAcc,
		LatestSessions: latestSessions,
	}, nil
}
