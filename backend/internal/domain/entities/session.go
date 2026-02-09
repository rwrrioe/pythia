package entities

import "time"

type Session struct {
	Id        int64
	Name      string `json:"name"`
	StartedAt time.Time
	EndedAt   time.Time `json:"ended_at"`
	Duration  time.Duration
	Status    string
	Language  int `json:"language"`
	Level     int `json:"level"`
	Accuracy  float64
}
