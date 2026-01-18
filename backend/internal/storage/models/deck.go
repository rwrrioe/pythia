package models

type Deck struct {
	Id        int `db:"id"`
	SessionId int `db:"session_id"`
}
