package ocr

import (
	"context"

	ocrclient "github.com/rwrrioe/pythia/backend/internal/clients/ocr/grpc"
)

type OCRService struct {
	Client *ocrclient.Client
}

func New(cl *ocrclient.Client) *OCRService {
	return &OCRService{Client: cl}
}

func (s *OCRService) ProcessImage(ctx context.Context, img []byte, lang string) ([]string, error) {
	text, err := s.Client.ProcessImage(ctx, img, lang)
	if err != nil {
		return nil, err
	}

	return text, nil
}
