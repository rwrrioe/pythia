package models

type User struct {
	Id          int    `db:"id"`
	Email       string `db:"email"`
	Name        string `db:"name"`
	Level       int    `db:"level_id"`
	Lang        int    `db:"lang_id"`
	WordsPerDay int    `db:"words_per_day"`
}
