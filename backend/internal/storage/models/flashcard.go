package models

import "github.com/google/uuid"

type FlashCard struct {
	Id     uuid.UUID `db:"id"`
	Word   string    `db:"word"`
	Transl string    `db:"transl"`
	Lang   int       `db:"lang_id"`
}
