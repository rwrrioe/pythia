package domain

import (
	"context"

	"github.com/rwrrioe/pythia/internal/domain/models"
)

type ImageProcessor interface {
	ProcessImage(ctx context.Context, imagePath string, lang string) (*models.OCRResult, error)
}

type WordTranslator interface {
	FindUnknown(ctx context.Context, text models.Text) ([]models.Word, error)
	TranslateWord(ctx context.Context, word string, fromLang string, toLang string) (*models.TranslatedWord, error)
}

type CardsBuilder interface {
	BuildCards(ctx context.Context, words []models.Word, transl []models.TranslatedWord) (*[]models.FlashCard, error)
}
