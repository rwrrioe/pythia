package domain

import (
	"context"

	"github.com/rwrrioe/pythia/internal/domain/entities"
	"github.com/rwrrioe/pythia/internal/domain/models"
)

type ImageProcessor interface {
	ProcessImage(ctx context.Context, imagePath string, lang string) (*models.OCRResult, error)
}

type WordTranslator interface {
	FindUnknownWords(ctx context.Context, text models.Text) ([]models.Word, error)
}

type CardsBuilder interface {
	BuildCards(ctx context.Context, words []models.Word, transl []models.TranslatedWord) (*[]models.FlashCard, error)
}

type TranslateProvider interface {
	FindUnknownWords(ctx context.Context, req models.AnalyzeRequest) ([]entities.UnknownWord, error)
}
