package database

import (
	"fmt"
	"os"

	"github.com/rwrrioe/pythia/internal/domain/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func DbConnect() (*gorm.DB, error) {
	host := os.Getenv("DB_HOST")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")
	port := os.Getenv("DB_PORT")

	if host == "" || user == "" || dbname == "" || port == "" {
		return nil, fmt.Errorf("missing DB env variables: host=%s user=%s dbname=%s port=%s", host, user, dbname, port)
	}

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		host, user, password, dbname, port,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return db, nil
}

func DbMigrate(db *gorm.DB) error {
	return db.AutoMigrate(&models.Word{}, &models.Translation{}, &models.FlashCard{})
}
