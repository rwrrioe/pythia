package usecase

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/rwrrioe/pythia/backend/internal/storage/postgresql"
	"github.com/rwrrioe/pythia/backend/internal/user/domain"
	"github.com/rwrrioe/pythia/backend/internal/user/ports"
)

type UserService struct {
	UserProvider       ports.UserProvider
	SessionProvider    ports.SessionProvider
	FlashCardsProvider ports.FlashCardProvider
	txm                *postgresql.TxManager
}

func (s *UserService) UserStats(
	ctx context.Context,
	uid uuid.UUID,
) (*domain.UserStats, error) {
	const op = "service.UserService.UserStats"

	usr, err := s.UserProvider.GetUser(ctx, uid)
	if err != nil {
		return nil, fmt.Errorf("%s:%w", op, err)
	}

	ss, err := s.SessionProvider.ListSessions(ctx, uid)
	if err != nil {
		return nil, fmt.Errorf("%s:%w", op, err)
	}

	flcards, err := s.FlashCardsProvider.ListFlashCards(ctx, uid)
	if err != nil {
		return nil, fmt.Errorf("%s:%w", op, err)
	}

	avgAcc := 0
	for _, s := range ss {
		avgAcc += int(s.Accuracy)
	}
	avgAcc = avgAcc / len(ss)

	return &domain.UserStats{
		WordsLearned: len(flcards),
		AvgAccuracy:  avgAcc,
		TotalSession: len(ss),
		Preferences: domain.UserLearningPreferences{
			Lang:      usr.Lang,
			Level:     usr.Level,
			DailyGoal: usr.WordsPerDay,
		},
	}, nil
}
