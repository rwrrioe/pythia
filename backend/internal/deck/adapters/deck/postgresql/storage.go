package postgresql

import decksqlc "github.com/rwrrioe/pythia/backend/internal/deck/adapters/postgresql/sqlc"

type Storage struct {
	q *decksqlc.Queries
}

func New(db decksqlc.DBTX) *Storage {
	return &Storage{
		q: decksqlc.New(db),
	}
}
