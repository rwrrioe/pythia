package ocr

import (
	"context"

	pb "github.com/rwrrioe/pythia/backend/shared/gen/go/ocr"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type ImageProcesser interface {
	ProcessImage(ctx context.Context, imageData []byte, lang string) ([]string, error)
}

type OCRProcessor struct{}

func NewOCRProcessor() *OCRProcessor {
	return &OCRProcessor{}
}

func (p *OCRProcessor) RecognizeText(ctx context.Context, imageData []byte, lang string) ([]string, error) {
	conn, err := grpc.NewClient(
		"localhost:50051",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)

	if err != nil {
		return nil, nil
	}
	defer conn.Close()

	client := pb.NewOCRServiceClient(conn)

	resp, err := client.Recognize(ctx, &pb.OCRRequest{
		ImageData: []byte(imageData),
		Lang:      lang,
	})
	if err != nil {
		return nil, err
	}

	return resp.Text, nil
}
