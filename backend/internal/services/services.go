package services

import (
	"context"

	cards "github.com/rwrrioe/pythia/backend/internal/services/cards_service"
	learn "github.com/rwrrioe/pythia/backend/internal/services/learn_service"
	"github.com/rwrrioe/pythia/backend/internal/services/ocr_service/ocr"
	translate "github.com/rwrrioe/pythia/backend/internal/services/translate_service"
)

type Services struct {
	Cards     *cards.FlashCardsService
	OCR       *ocr.OCRProcesser
	Translate *translate.TranslateService
	Learn     *learn.LearnService
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

	learn := learn.NewLearnService(4)
	return &Services{
		Cards:     cards,
		OCR:       ocr,
		Translate: transl,
		Learn:     learn,
	}, nil
}
