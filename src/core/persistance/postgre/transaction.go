package postgre

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/Shitcode-Swamp/go-train-framework/src/core/persistance/postgre/internal/pgsql"
	"github.com/Shitcode-Swamp/go-train-framework/src/core/persistance/store"
)

type Database struct {
	pool *pgxpool.Pool
}

func NewDatabase(pool *pgxpool.Pool) store.Transactor {
	return &Database{pool: pool}
}

func (d *Database) WithTransaction(ctx context.Context, fn func(ctx context.Context) error, options ...store.TransactionOptions) error {
	tx, err := d.pool.BeginTx(ctx, pgx.TxOptions{IsoLevel: isoLevel(options)})
	if err != nil {
		return pgsql.TranslateError(err)
	}
	defer tx.Rollback(ctx)

	if err := fn(pgsql.ContextWithTx(ctx, tx)); err != nil {
		return err
	}
	return pgsql.TranslateError(tx.Commit(ctx))
}

func isoLevel(options []store.TransactionOptions) pgx.TxIsoLevel {
	if len(options) == 0 {
		return pgx.ReadCommitted
	}
	switch options[0].IsolationLevel {
	case store.IsolationSerializable:
		return pgx.Serializable
	case store.IsolationRepeatableRead:
		return pgx.RepeatableRead
	default:
		return pgx.ReadCommitted
	}
}

var _ store.Transactor = (*Database)(nil)
