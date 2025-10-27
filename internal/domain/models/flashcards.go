package models

import "time"

type FlashCard struct {
	CardID        uint `gorm:"column:id;primaryKey;autoIncrement; uniqueIndex:index_flashcardid"`
	WordID        uint `gorm:"not null"`
	TranslationID uint `gorm:"not null"`
	BatchId       int
	CreatedAt     time.Time `gorm:"type:uuid;index"`
}

func (f FlashCard) TableName() string {
	return "flashcards"
}
