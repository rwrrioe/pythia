package domain

//интерфейсы поближе к коду лучше
import (
	"context"

	"github.com/rwrrioe/pythia/backend/internal/domain/entities"
	"github.com/rwrrioe/pythia/backend/internal/domain/models"
)

type ImageProcesser interface {
	ProcessImage(ctx context.Context, imageData []byte, lang string) ([]string, error)
}

type CardsBuilder interface {
	BuildCards(ctx context.Context, words []models.Word, transl []models.TranslatedWord) (*[]models.FlashCard, error)
}

type TranslateProvider interface {
	FindUnknownWords(ctx context.Context, req models.AnalyzeRequest) ([]entities.UnknownWord, error)
	WriteExamples(ctx context.Context, words []entities.UnknownWord, req models.AnalyzeRequest) (*[]entities.Example, error)
}

type TestsProvider interface {
	QuizTest(ctx context.Context, words *[]entities.UnknownWord) *[]entities.QuizQuestion
}
