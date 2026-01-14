package postgresql

import (
	"context"

	"github.com/google/uuid"
	"github.com/rwrrioe/pythia/backend/internal/domain/models"
	"gorm.io/gorm"
)

type FlashCardRepository struct {
	db *gorm.DB
}

func NewFlashCardRepo(db *gorm.DB) *FlashCardRepository {
	return &FlashCardRepository{
		db: db,
	}
}

func (fl *FlashCardRepository) GetById(ctx context.Context, cardId int) (*models.FlashCard, error) {
	var response models.FlashCard
	result := fl.db.WithContext(ctx).First(&response, cardId)
	if result.Error != nil {
		return nil, result.Error
	}

	return &response, nil
}

func (fl *FlashCardRepository) Save(ctx context.Context, card models.FlashCard) error {
	result := fl.db.WithContext(ctx).Create(&card)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (fl *FlashCardRepository) ListByUUID(ctx context.Context, batchID uuid.UUID) ([]models.FlashCard, error) {
	var flashcards []models.FlashCard
	result := fl.db.WithContext(ctx).Where("batch_id=?", batchID).Find(&flashcards).Order("created_at ASC")
	if result.Error != nil {
		return nil, result.Error
	}

	return flashcards, nil
}
