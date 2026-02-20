package sso_grpc_client

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/google/uuid"
	ssov2 "github.com/rwrrioe/sso_protos/v2/gen/go/sso/sso"
	"google.golang.org/grpc"
)

type Client struct {
	api ssov2.AuthClient
	log *slog.Logger
}

func New(
	cc *grpc.ClientConn,
	log *slog.Logger,
) *Client {
	grpcClient := ssov2.NewAuthClient(cc)

	return &Client{
		api: grpcClient,
		log: log,
	}
}

func (c *Client) IsAdmin(ctx context.Context, userID uuid.UUID) (bool, error) {
	const op = "sso_grpc.IsAdmin"

	resp, err := c.api.IsAdmin(ctx, &ssov2.IsAdminRequest{
		UserId: userID.String(),
	})
	if err != nil {
		return false, fmt.Errorf("%s: %w", op, err)
	}

	return resp.IsAdmin, nil
}

func (c *Client) Login(ctx context.Context, email string, passwd string, appId int32) (string, error) {
	const op = "sso_grpc.Login"

	resp, err := c.api.Login(ctx, &ssov2.LoginRequest{
		Email:    email,
		Password: passwd,
		AppId:    appId,
	})

	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return resp.Token, nil
}

func (c *Client) Register(ctx context.Context, email string, passwd string) (uuid.UUID, error) {
	const op = "sso_grpc.Register"

	resp, err := c.api.Register(ctx, &ssov2.RegisterRequest{
		Email:    email,
		Password: passwd,
	})

	if err != nil {
		return uuid.Nil, fmt.Errorf("%s: %w", op, err)
	}

	uid, err := uuid.Parse(resp.UserId)
	if err != nil {
		return uuid.Nil, fmt.Errorf("%s:%w", op, "failed to parse uid to uuid")
	}

	return uid, nil
}
