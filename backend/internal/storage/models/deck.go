package models

type Deck struct {
	UserId    int   `db:"user_id"`
	Id        int   `db:"id"`
	SessionId int64 `db:"session_id"`
}
