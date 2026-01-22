package entities

type FlashCardDTO struct {
	Word        string `json:"word"`
	Translation string `json:"translation"`
	Lang        string `json:"language"`
}

type FlashCard struct {
	Id     int
	Word   string
	Transl string
	Desc   string
	Lang   int
}
