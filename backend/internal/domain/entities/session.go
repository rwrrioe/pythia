package entities

import (
	"time"

	"github.com/google/uuid"
)

type Session struct {
	Id        uuid.UUID
	Name      string `json:"name"`
	StartedAt time.Time
	EndedAt   time.Time `json:"ended_at"`
	Duration  time.Duration
	Status    string
	Language  int `json:"language"`
	Level     int `json:"level"`
	Accuracy  float64
}
