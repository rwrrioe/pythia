package entities

import "time"

type Session struct {
	Id        int64
	UserId    int64
	Name      string
	StartedAt time.Time
	EndedAt   time.Time
	Duration  time.Duration
	Status    string
	Language  int
	Level     int
	Accuracy  float64
}
