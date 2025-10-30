package models

import "time"

type Word struct {
	Id          uint   `gorm:"primaryKey;autoIncrement; uniqueIndex:index_wordid"`
	Text        string `gorm:"not null"`
	Language    string `gorm:"default:en"`
	CreatedAt   time.Time
	Translation Translation `gorm:"foreignKey:WordID"`
}

type Translation struct {
	Id             uint   `gorm:"primaryKey;autoIncrement; uniqueIndex:index_translationid"`
	WordID         uint   `gorm:"not null;unique"`
	TranslatedWord string `gorm:"not null"`
	Language       string `gorm:"not null"`
	CreatedAt      time.Time
}
