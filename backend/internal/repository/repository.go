package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/rwrrioe/pythia/backend/internal/domain/models"
)

type FlashCardRepo interface {
	GetById(ctx context.Context, cardID int) (*models.FlashCard, error)
	Save(ctx context.Context, card models.FlashCard) error
	ListByUUID(ctx context.Context, batchID uuid.UUID) ([]models.FlashCard, error)
}

type TranslationRepo interface {
	GetById(ctx context.Context, translId int) (*models.Translation, error)
	Save(ctx context.Context, transl models.Translation) error
}

type WordRepo interface {
	GetById(ctx context.Context, wordId int) (*models.Word, error)
	Save(ctx context.Context, word models.Word) error
}
