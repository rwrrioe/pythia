package models

import "time"

type Session struct {
	Id        int       `db:"id"`
	Name      string    `db:"name"`
	StartedAt time.Time `db:"started_at"`
	EndedAt   time.Time `db:"ended_at"`
	Accuracy  float64   `db:"accuracy"`
}
