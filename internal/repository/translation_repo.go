package repository

import (
	"context"

	"github.com/rwrrioe/pythia/internal/domain/models"
	"gorm.io/gorm"
)

type TranslationRepo interface {
	GetById(ctx context.Context, translId int) (*models.Translation, error)
	Save(ctx context.Context, transl models.Translation) error
}

type TranslationRepository struct {
	db *gorm.DB
}

func NewTranslationRepo(db *gorm.DB) *TranslationRepository {
	return &TranslationRepository{
		db: db,
	}
}

func (tr *TranslationRepository) GetById(ctx context.Context, translId int) (*models.Translation, error) {
	var response models.Translation
	result := tr.db.WithContext(ctx).First(&response, translId)
	if result.Error != nil {
		return nil, result.Error
	}

	return &response, nil
}

func (tr *TranslationRepository) Save(ctx context.Context, transl models.Translation) error {
	result := tr.db.WithContext(ctx).Create(&transl)
	if result.Error != nil {
		return result.Error
	}

	return nil
}
