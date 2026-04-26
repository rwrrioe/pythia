package entities

import "github.com/google/uuid"

type Deck struct {
	Id         uuid.UUID
	SessionId  uuid.UUID
	Flashcards []FlashCard
}
