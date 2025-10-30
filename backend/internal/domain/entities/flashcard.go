package entities

type FlashCardDTO struct {
	Word        string `json:"word"`
	Translation string `json:"translation"`
	Lang        string `json:"language"`
}
