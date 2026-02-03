package models

type User struct {
	Id          int64  `db:"id"`
	Email       string `db:"email"`
	Name        string `db:"name"`
	Level       string `db:"level"`
	Lang        string `db:"language"`
	WordsPerDay int    `db:"words_per_day"`
}
