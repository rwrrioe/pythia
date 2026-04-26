package authz

import (
	"context"

	"github.com/google/uuid"
)

type SessionProvider interface {
	GetSession(
		ctx context.Context,
		sessionId uuid.UUID,
	) (uuid.UUID, bool, error)
}

type SessionAccessPolicy struct {
	session SessionProvider
}
