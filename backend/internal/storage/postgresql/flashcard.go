package postgresql

import (
	"github.com/jackc/pgx/v5"
)

type FlashCardStorage struct {
	conn *pgx.Conn
}

func NewFlashcardRepo(conn *pgx.Conn) *FlashCardStorage {
	return &FlashCardStorage{
		conn: conn,
	}
}
