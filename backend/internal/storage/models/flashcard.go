package models

type FlashCard struct {
	Id     int    `db:"id"`
	Word   string `db:"word"`
	Transl string `db:"transl"`
	Lang   int    `db:"lang_id"`
}
