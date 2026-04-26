package entities

import "github.com/google/uuid"

type FlashCardDTO struct {
	Word        string `json:"word"`
	Translation string `json:"translation"`
	Lang        string `json:"language"`
}

type FlashCard struct {
	Id     uuid.UUID
	Word   string
	Transl string
	Desc   string
	Lang   int
}
