package postgresql

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type TxManager struct {
	Pool *pgxpool.Pool
}

func NewTxManager(pool *pgxpool.Pool) *TxManager {
	return &TxManager{Pool: pool}
}

func (m *TxManager) WithTx(ctx context.Context, fn func(tg pgx.Tx) error) error {
	const op = "postgresql.TxManager.WithTx"
	tx, err := m.Pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return fmt.Errorf("%s:%w", op, err)
	}

	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback(ctx)
		}
		if err != nil {
			_ = tx.Rollback(ctx)
		}

		err = tx.Commit(ctx)
	}()

	return fn(tx)
}
