package service

import (
	"context"
	"fmt"

	"github.com/rwrrioe/pythia/backend/internal/auth/authn"
	"github.com/rwrrioe/pythia/backend/internal/domain/entities"
	service "github.com/rwrrioe/pythia/backend/internal/services/errors"
	"github.com/rwrrioe/pythia/backend/internal/storage/postgresql"
)

type UserProvider interface {
	GetUser(ctx context.Context) (*entities.User, error)
}

type UserService struct {
	User       UserProvider
	Session    SessionProvider
	FlashCards FlashCardProvider
	txm        *postgresql.TxManager
}

func (s *UserService) UserStats(ctx context.Context) (*entities.UserStats, error) {
	const op = "service.UserService.UserStats"

	uid, ok := authn.UIDFromContext(ctx)
	if !ok {
		return nil, service.ErrUnauthorized
	}

	usr, err := s.User.GetUser(ctx)
	if err != nil {
		return nil, fmt.Errorf("%s:%w", op, err)
	}

	ss, err := s.Session.ListSessions(ctx, s.txm.Pool, uid)
	if err != nil {
		return nil, fmt.Errorf("%s:%w", op, err)
	}

	flcards, err := s.FlashCards.List(ctx, s.txm.Pool, uid)
	if err != nil {
		return nil, fmt.Errorf("%s:%w", op, err)
	}

	avgAcc := 0
	for _, s := range ss {
		avgAcc += int(s.Accuracy)
	}
	avgAcc = avgAcc / len(ss)

	return &entities.UserStats{
		WordsLearned: len(flcards),
		AvgAccuracy:  avgAcc,
		TotalSession: len(ss),
		Preferences: entities.UserLearningPreferences{
			Lang:      usr.Lang,
			Level:     usr.Level,
			DailyGoal: usr.WordsPerDay,
		},
	}, nil
}
