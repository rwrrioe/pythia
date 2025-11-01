package ocr

import (
	"context"

	pb "github.com/rwrrioe/pythia/shared/gen/go/ocr"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type ImageProcesser interface {
	ProcessImage(ctx context.Context, imageData []byte, lang string) ([]string, error)
}

type OCRProcessor struct {
	client pb.OCRServiceClient
	conn   *grpc.ClientConn
}

func NewOCRProcessor(add string) (*OCRProcessor, error) {
	conn, err := grpc.NewClient(
		add,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, err
	}
	client := pb.NewOCRServiceClient(conn)
	return &OCRProcessor{
		client: client,
		conn:   conn,
	}, nil
}

func (p *OCRProcessor) Close() error {
	return p.conn.Close()
}

func (p *OCRProcessor) RecognizeText(ctx context.Context, imageData []byte, lang string) ([]string, error) {
	resp, err := p.client.Recognize(ctx, &pb.OCRRequest{
		ImageData: []byte(imageData),
		Lang:      lang,
	})
	if err != nil {
		return nil, err
	}

	return resp.Text, nil
}
