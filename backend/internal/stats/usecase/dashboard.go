package usecase

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/rwrrioe/pythia/backend/internal/stats/domain"
	"github.com/rwrrioe/pythia/backend/internal/stats/domain/errors"
	"github.com/rwrrioe/pythia/backend/internal/stats/ports"
	"github.com/rwrrioe/pythia/backend/internal/storage/postgresql"
)

type Stats struct {
	sessionProvider   ports.SessionProvider
	userProvider      ports.UserProvider
	flashCardProvider ports.FlashCardProvider
}

func NewStatsService(
	ss ports.SessionProvider,
	fl ports.FlashCardProvider,
) *Stats {
	return &Stats{
		sessionProvider:   ss,
		flashCardProvider: fl,
	}
}

func (s *Stats) Dashboard(
	ctx context.Context,
	uid uuid.UUID,
) (*domain.Dashboard, error) {
	const op = "service.Stats.Dashboard"

	latestSessions, err := s.sessionProvider.ListLatest(ctx, uid)
	if err != nil {
		if errors.Is(err, postgresql.ErrUserNotFound) {
			return nil, fmt.Errorf("%s:%w", op, domainerrors.ErrUnauthorized)
		}

		return nil, fmt.Errorf("%s:%w", op, err)
	}

	words, err := s.flashCardProvider.ListFlashCards(ctx, uid)
	if err != nil {
		return nil, fmt.Errorf("%s:%w", op, err)
	}

	sessions, err := s.sessionProvider.ListSessions(ctx, uid)
	if err != nil {
		if errors.Is(err, postgresql.ErrUserNotFound) {
			return nil, fmt.Errorf("%s:%w", op, domainerrors.ErrUnauthorized)
		}

		return nil, fmt.Errorf("%s:%w", op, err)
	}

	if len(sessions) == 0 {
		return &domain.Dashboard{
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

	return &domain.Dashboard{
		Streak:         0,
		WordsLearned:   len(words),
		Accuracy:       avgAcc,
		LatestSessions: latestSessions,
	}, nil
}
