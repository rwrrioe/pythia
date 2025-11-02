package ocr

import (
	"context"

	pb "github.com/rwrrioe/pythia/shared/gen/go/ocr"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type OCRProcesser struct {
	client pb.OCRServiceClient
	conn   *grpc.ClientConn
}

func NewOCRProcessor(add string) (OCRProcesser, error) {
	conn, err := grpc.NewClient(
		add,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return OCRProcesser{}, err
	}
	client := pb.NewOCRServiceClient(conn)
	return OCRProcesser{
		client: client,
		conn:   conn,
	}, nil
}

func (p *OCRProcesser) Close() error {
	return p.conn.Close()
}

func (p *OCRProcesser) ProcessImage(ctx context.Context, imageData []byte, lang string) ([]string, error) {
	resp, err := p.client.Recognize(ctx, &pb.OCRRequest{
		ImageData: []byte(imageData),
		Lang:      lang,
	})
	if err != nil {
		return nil, err
	}

	return resp.Text, nil
}
