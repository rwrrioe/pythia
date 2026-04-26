package usecase

import (
	"context"

	"github.com/rwrrioe/pythia/backend/ocr/ports"
)

type OCR struct {
	OCRProvider ports.OCRProvider
}

func NewOCRService(ocr ports.OCRProvider) *OCR {
	return &OCR{OCRProvider: ocr}
}

func (s *OCR) ProcessImage(ctx context.Context, img []byte, lang string) ([]string, error) {
	text, err := s.OCRProvider.ProcessImage(ctx, img, lang)
	if err != nil {
		return nil, err
	}

	return text, nil
}
