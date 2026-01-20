package service

import (
	"context"

	ocrclient "github.com/rwrrioe/pythia/backend/internal/clients/ocr/grpc"
)

type Services struct {
	CardService   *FlashCardsService
	OcrService    *OCRService
	TranslService *TranslateService
	LearnService  *LearnService
}

func New(ctx context.Context, ocrClient *ocrclient.Client, aiModel string) (*Services, error) {
	cardService := NewCardsService()
	// ocrService := New(ocrClient)

	translService, err := NewTranslateService(ctx, aiModel)
	if err != nil {
		return nil, err
	}

	learnService := NewLearnService(4)
	return &Services{
		CardService: cardService,
		// OcrService:    ocrService,
		TranslService: translService,
		LearnService:  learnService,
	}, nil
}
