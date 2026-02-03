package authn

import (
	"context"
	"errors"
	"fmt"

	sso_grpc_client "github.com/rwrrioe/pythia/backend/internal/clients/sso/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	// Для HTTP handlers: 401
	ErrInvalidCredentials = errors.New("invalid credentials")

	// Для 409
	ErrUserAlreadyExists = errors.New("user already exists")

	// Для 503/500
	ErrSSOUnavailable = errors.New("sso unavailable")
)

type SSOService interface {
	Login(ctx context.Context, email, password string) (string, error)
	Register(ctx context.Context, email, password string) (int64, error)
	IsAdmin(ctx context.Context, userID int64) (bool, error)
}

type service struct {
	sso   *sso_grpc_client.Client
	appID int32
}

func NewSSO(sso *sso_grpc_client.Client, appID int32) SSOService {
	return &service{
		sso:   sso,
		appID: appID}
}

func (s *service) Login(ctx context.Context, email, password string) (string, error) {
	token, err := s.sso.Login(ctx, email, password, s.appID)
	if err == nil {
		return token, nil
	}
	return "", mapSSOErr("auth.Login", err)
}

func (s *service) Register(ctx context.Context, email, password string) (int64, error) {
	uid, err := s.sso.Register(ctx, email, password)
	if err == nil {
		return uid, nil
	}
	return 0, mapSSOErr("auth.Register", err)
}

func (s *service) IsAdmin(ctx context.Context, userID int64) (bool, error) {
	ok, err := s.sso.IsAdmin(ctx, userID)
	if err == nil {
		return ok, nil
	}
	return false, mapSSOErr("auth.IsAdmin", err)
}

func mapSSOErr(op string, err error) error {
	// unwrap grpc status if present
	st, ok := status.FromError(err)
	if !ok {
		return fmt.Errorf("%s: %w", op, ErrSSOUnavailable)
	}

	switch st.Code() {
	case codes.Unauthenticated:
		return fmt.Errorf("%s: %w", op, ErrInvalidCredentials)

	case codes.AlreadyExists:
		return fmt.Errorf("%s: %w", op, ErrUserAlreadyExists)

	case codes.Unavailable, codes.DeadlineExceeded:
		return fmt.Errorf("%s: %w", op, ErrSSOUnavailable)

	default:
		return fmt.Errorf("%s: %w", op, err)
	}
}
