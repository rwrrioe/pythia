package sso_grpc_client

import (
	"context"
	"log/slog"

	ssov1 "github.com/rwrrioe/sso_protos/gen/go/sso"
	"google.golang.org/grpc"
)

type Client struct {
	api ssov1.AuthClient
	log *slog.Logger
}

func New(
	ctx context.Context,
	cc *grpc.ClientConn,
	log *slog.Logger,
) *Client {
	grpcClient := ssov1.NewAuthClient(cc)

	return &Client{
		api: grpcClient,
		log: log,
	}
}
