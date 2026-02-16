package authz

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/google/uuid"
	taskstorage "github.com/rwrrioe/pythia/backend/internal/storage/redis/task_storage"
)

var (
	ErrForbidden       = fmt.Errorf("access forbidden")
	ErrSessionNotFound = fmt.Errorf("session not found")
)

type AuthorizeService interface {
	CanAccessSession(ctx context.Context, uid int64, sessionId uuid.UUID) error
}

type authorizer struct {
	redis *taskstorage.RedisStorage
	log   *slog.Logger
}

func NewAuthorizer(redis *taskstorage.RedisStorage, log *slog.Logger) AuthorizeService {
	return &authorizer{
		redis: redis,
		log:   log,
	}
}

func (a *authorizer) CanAccessSession(ctx context.Context, uid int64, sessionId uuid.UUID) error {
	const op = "authz.authorizer.CanAccessSession"

	ss, ok, err := a.redis.GetSession(ctx, sessionId)
	if err != nil {
		return fmt.Errorf("%s:%w", op, err)
	} else if !ok {
		return fmt.Errorf("%s:%w", op, ErrSessionNotFound)
	}

	if uid != ss.UserId {
		val := fmt.Sprintf("have %d want %d", uid, ss.UserId)
		a.log.Warn("access forbidden", slog.String("user ids", val))
		return fmt.Errorf("%s:%w", op, ErrForbidden)
	}

	return nil
}
