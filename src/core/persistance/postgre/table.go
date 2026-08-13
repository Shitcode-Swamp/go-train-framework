package postgre

import (
	"context"
	"reflect"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/Shitcode-Swamp/go-train-framework/src/core/persistance/postgre/internal/pgsql"
	"github.com/Shitcode-Swamp/go-train-framework/src/core/persistance/store"
)

type querier interface {
	Exec(ctx context.Context, sql string, args ...any) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
}

type Table[T any] struct {
	pool    *pgxpool.Pool
	name    string
	mapping pgsql.TableMapping
}

func New[T any](pool *pgxpool.Pool, table string) store.Store[T] {
	var zero T
	return &Table[T]{
		pool:    pool,
		name:    table,
		mapping: pgsql.MappingFor(reflect.TypeOf(zero)),
	}
}

var _ store.Store[struct{}] = (*Table[struct{}])(nil)

func (s *Table[T]) columnNames() []string {
	cols := make([]string, len(s.mapping.Columns))
	for i, c := range s.mapping.Columns {
		cols[i] = c.Column
	}
	return cols
}

func (s *Table[T]) withTx(ctx context.Context, fn func(querier) error) error {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return pgsql.TranslateError(err)
	}
	defer tx.Rollback(ctx)

	if err := fn(tx); err != nil {
		return err
	}
	return pgsql.TranslateError(tx.Commit(ctx))
}
