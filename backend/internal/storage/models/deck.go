package models

import "github.com/google/uuid"

type Deck struct {
	UserId    uuid.UUID `db:"user_id"`
	Id        uuid.UUID `db:"id"`
	SessionId uuid.UUID `db:"session_id"`
}
