package models

import "github.com/google/uuid"

type Text struct {
	Text []byte
	UUID uuid.UUID
}
