package models

import "github.com/google/uuid"

type User struct {
	Id          uuid.UUID `db:"id"`
	Email       string    `db:"email"`
	Name        string    `db:"name"`
	Level       string    `db:"level"`
	Lang        string    `db:"language"`
	WordsPerDay int       `db:"words_per_day"`
}
