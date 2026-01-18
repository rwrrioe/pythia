package service

import (
	"context"

	ocrclient "github.com/rwrrioe/pythia/backend/internal/clients/ocr/grpc"
	cards "github.com/rwrrioe/pythia/backend/internal/services/cards_service"
	learn "github.com/rwrrioe/pythia/backend/internal/services/learn_service"
	ocr "github.com/rwrrioe/pythia/backend/internal/services/ocr_service/ocr"
	translate "github.com/rwrrioe/pythia/backend/internal/services/translate_service"
)

type Services struct {
	CardService   *cards.FlashCardsService
	OcrService    *ocr.OCRService
	TranslService *translate.TranslateService
	LearnService  *learn.LearnService
}

func New(ctx context.Context, ocrClient *ocrclient.Client, aiModel string) (*Services, error) {
	cardService := cards.NewCardsService()
	ocrService := ocr.New(ocrClient)

	translService, err := translate.NewTranslateService(ctx, aiModel)
	if err != nil {
		return nil, err
	}

	learnService := learn.NewLearnService(4)
	return &Services{
		CardService:   cardService,
		OcrService:    ocrService,
		TranslService: translService,
		LearnService:  learnService,
	}, nil
}
