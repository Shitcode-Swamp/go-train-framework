package postgre

import (
	"context"
	"fmt"
	"reflect"
	"strings"

	"github.com/Shitcode-Swamp/go-train-framework/src/core/store/contract"
	"github.com/Shitcode-Swamp/go-train-framework/src/core/store/filters"
	"github.com/Shitcode-Swamp/go-train-framework/src/core/store/postgre/internal/pgsql"
)

func (s *Table[T]) Update(ctx context.Context, filter filters.Filtered[T], update T, options ...contract.UpdateOptions) error {
	upsert := false
	if len(options) > 0 && options[0].Upsert != nil {
		upsert = *options[0].Upsert
	}
	return s.update(ctx, s.pool, filter, update, upsert)
}

func (s *Table[T]) update(ctx context.Context, q querier, filter filters.Filtered[T], update T, upsert bool) error {
	setClauses, args := s.updateSetClause(update)
	where, whereArgs := pgsql.BuildWhere(filter, len(args))
	args = append(args, whereArgs...)

	query := fmt.Sprintf("UPDATE %s SET %s WHERE %s", s.name, strings.Join(setClauses, ", "), where)
	tag, err := q.Exec(ctx, query, args...)
	if err != nil {
		return pgsql.TranslateError(err)
	}

	if tag.RowsAffected() == 0 && upsert {
		_, err := s.insert(ctx, q, update)
		return pgsql.TranslateError(err)
	}
	return nil
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

func (s *Table[T]) BulkUpdate(ctx context.Context, ops []contract.BulkUpdateOp[T]) error {
	if len(ops) == 0 {
		return nil
	}

	return s.withTx(ctx, func(q querier) error {
		for _, op := range ops {
			if err := s.update(ctx, q, op.Filter, op.Update, op.Upsert); err != nil {
				return err
			}
		}
		return nil
	})
}
