package ocr_grpc_client

import (
	"context"
	"log/slog"

	ocrv1 "github.com/rwrrioe/pythia_protos/gen/go/ocr"
	"google.golang.org/grpc"
)

type Client struct {
	api ocrv1.OCRServiceClient
	log *slog.Logger
}

func New(
	ctx context.Context,
	cc *grpc.ClientConn,
	log *slog.Logger,
) *Client {
	grpcClient := ocrv1.NewOCRServiceClient(cc)

	return &Client{
		api: grpcClient,
		log: log,
	}
}

func (c *Client) ProcessImage(ctx context.Context, imageData []byte, lang string) ([]string, error) {
	resp, err := c.api.Recognize(ctx, &ocrv1.OCRRequest{
		ImageData: []byte(imageData),
		Lang:      lang,
	})
	if err != nil {
		return nil, err
	}

	return resp.Text, nil
}
