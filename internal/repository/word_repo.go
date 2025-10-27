package repository

import (
	"context"

	"github.com/rwrrioe/pythia/internal/domain/models"
	"gorm.io/gorm"
)

type WordRepo interface {
	GetById(ctx context.Context, wordId int) (*models.Word, error)
	Save(ctx context.Context, word models.Word) error
}

type WordRepository struct {
	db *gorm.DB
}

func NewWordRepo(db *gorm.DB) *WordRepository {
	return &WordRepository{
		db: db,
	}
}

func (wr *WordRepository) GetById(ctx context.Context, wordId int) (*models.Word, error) {
	var response models.Word
	result := wr.db.WithContext(ctx).First(&response, wordId)
	if result.Error != nil {
		return nil, result.Error
	}

	return &response, nil
}

func (wr *WordRepository) Save(ctx context.Context, word models.Word) error {
	result := wr.db.WithContext(ctx).Create(&word)
	if result.Error != nil {
		return result.Error
	}

	return nil
}
