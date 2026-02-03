package sso_grpc_client

import (
	"context"
	"fmt"
	"log/slog"

	ssov1 "github.com/rwrrioe/sso_protos/gen/go/sso"
	"google.golang.org/grpc"
)

type Client struct {
	api ssov1.AuthClient
	log *slog.Logger
}

func New(
	cc *grpc.ClientConn,
	log *slog.Logger,
) *Client {
	grpcClient := ssov1.NewAuthClient(cc)

	return &Client{
		api: grpcClient,
		log: log,
	}
}

func (c *Client) IsAdmin(ctx context.Context, userID int64) (bool, error) {
	const op = "sso_grpc.IsAdmin"

	resp, err := c.api.IsAdmin(ctx, &ssov1.IsAdminRequest{
		UserId: userID,
	})
	if err != nil {
		return false, fmt.Errorf("%s: %w", op, err)
	}

	return resp.IsAdmin, nil
}

func (c *Client) Login(ctx context.Context, email string, passwd string, appId int32) (string, error) {
	const op = "sso_grpc.Login"

	resp, err := c.api.Login(ctx, &ssov1.LoginRequest{
		Email:    email,
		Password: passwd,
		AppId:    appId,
	})

	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return resp.Token, nil
}

func (c *Client) Register(ctx context.Context, email string, passwd string) (int64, error) {
	const op = "sso_grpc.Login"

	resp, err := c.api.Register(ctx, &ssov1.RegisterRequest{
		Email:    email,
		Password: passwd,
	})

	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return resp.UserId, nil
}
