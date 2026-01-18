package models

type FlashCard struct {
	Id     int    `db:"id"`
	Word   string `db:"word"`
	Transl string `db:"transl"`
	Desc   string `db:"description"`
	Lang   string
}
