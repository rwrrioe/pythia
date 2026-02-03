package models

import "time"

type Session struct {
	Id        int64     `db:"id"`
	UserId    int64     `db:"user_id"`
	Name      string    `db:"name"`
	Lang      int       `db:"lang_id"`
	Status    string    `db:"status"`
	StartedAt time.Time `db:"started_at"`
	EndedAt   time.Time `db:"ended_at"`
	Accuracy  float64   `db:"accuracy"`
}
