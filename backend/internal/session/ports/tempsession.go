package ports

import (
	"context"

	"github.com/google/uuid"
	"github.com/rwrrioe/pythia/backend/internal/storage/redis/task_storage"
)

type RedisProvider interface {
	GetBySession(ctx context.Context, sessionId uuid.UUID) ([]redis_storage.TaskDTO, bool, error)
	SaveSession(ctx context.Context, ss redis_storage.SessionDTO) error
	UpdateSession(ctx context.Context, sessionId uuid.UUID, update func(s *redis_storage.SessionDTO)) (bool, error)
}
