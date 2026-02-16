package models

import "github.com/google/uuid"

type Deck struct {
	UserId    int       `db:"user_id"`
	Id        uuid.UUID `db:"id"`
	SessionId int64     `db:"session_id"`
}
