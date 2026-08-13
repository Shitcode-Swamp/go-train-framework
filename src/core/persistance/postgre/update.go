package postgre

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/jackc/pgx/v5"

	"github.com/Shitcode-Swamp/go-train-framework/src/core/persistance/filters"
	"github.com/Shitcode-Swamp/go-train-framework/src/core/persistance/postgre/internal/pgsql"
	"github.com/Shitcode-Swamp/go-train-framework/src/core/persistance/store"
)

func (s *Table[T]) Update(ctx context.Context, filter filters.Filtered[T], update T, options ...store.UpdateOptions) (any, error) {
	upsert := false
	if len(options) > 0 && options[0].Upsert != nil {
		upsert = *options[0].Upsert
	}
	return s.update(ctx, s.pool, filter, update, upsert)
}

func (s *Table[T]) update(ctx context.Context, q querier, filter filters.Filtered[T], update T, upsert bool) (any, error) {
	setClauses, args := s.updateSetClause(update)
	where, whereArgs := pgsql.BuildWhere(filter, len(args))
	args = append(args, whereArgs...)

	query := fmt.Sprintf("UPDATE %s SET %s WHERE %s", s.name, strings.Join(setClauses, ", "), where)

	if s.mapping.PKIndex < 0 {
		tag, err := q.Exec(ctx, query, args...)
		if err != nil {
			return nil, pgsql.TranslateError(err)
		}
		if tag.RowsAffected() == 0 && upsert {
			return s.insert(ctx, q, update)
		}
		return nil, nil
	}

	id, err := s.updateReturningID(ctx, q, query, args)
	if err != nil {
		return nil, pgsql.TranslateError(err)
	}
	if id == nil && upsert {
		return s.insert(ctx, q, update)
	}
	return id, nil
}

func (s *Table[T]) updateReturningID(ctx context.Context, q querier, query string, args []any) (any, error) {
	query += fmt.Sprintf(" RETURNING %s", s.mapping.Columns[s.mapping.PKIndex].Column)
	var id any
	if err := q.QueryRow(ctx, query, args...).Scan(&id); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return id, nil
}

func (s *Table[T]) updateSetClause(update T) (clauses []string, args []any) {
	v := reflect.ValueOf(update)
	for _, c := range s.mapping.Columns {
		if c.PK {
			continue
		}
		args = append(args, v.Field(c.Index).Interface())
		clauses = append(clauses, fmt.Sprintf("%s = $%d", c.Column, len(args)))
	}
	return clauses, args
}

func (s *Table[T]) BulkUpdate(ctx context.Context, ops []store.BulkUpdateOp[T]) ([]any, error) {
	if len(ops) == 0 {
		return nil, nil
	}

	ids := make([]any, len(ops))
	err := s.withTx(ctx, func(q querier) error {
		for i, op := range ops {
			id, err := s.update(ctx, q, op.Filter, op.Update, op.Upsert)
			if err != nil {
				return err
			}
			ids[i] = id
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return ids, nil
}
