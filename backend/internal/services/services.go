package services

import (
	"context"

	cards "github.com/rwrrioe/pythia/backend/internal/services/cards_service"
	"github.com/rwrrioe/pythia/backend/internal/services/ocr_service/ocr"
	translate "github.com/rwrrioe/pythia/backend/internal/services/translate_service"
)

type Services struct {
	Cards     *cards.CardsService
	OCR       *ocr.OCRProcesser
	Translate *translate.TranslateService
}

func New(ctx context.Context, AImodel string, grpcadd string) (*Services, error) {
	cards := cards.NewCardsService()
	ocr, err := ocr.NewOCRProcessor(grpcadd)
	if err != nil {
		return nil, err
	}
	transl, err := translate.NewTranslateService(ctx, AImodel)
	if err != nil {
		return nil, err
	}

	return &Services{
		Cards:     cards,
		OCR:       ocr,
		Translate: transl,
	}, nil
}
