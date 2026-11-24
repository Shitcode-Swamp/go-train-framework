package postgre

import (
	"context"
	"fmt"
	"strings"

	"github.com/Shitcode-Swamp/go-train-framework/v3/src/core/persistance/filters"
	"github.com/Shitcode-Swamp/go-train-framework/v3/src/core/persistance/postgre/internal/pgsql"
)

func (s *Table[T]) selectQuery(filter filters.Filtered[T], limit int) (string, []any) {
	where, args := pgsql.BuildWhere(filter, 0)
	query := fmt.Sprintf("SELECT %s FROM %s WHERE %s", strings.Join(s.columnNames(), ", "), s.name, where)
	if limit > 0 {
		query += fmt.Sprintf(" LIMIT %d", limit)
	}
	return query, args
}

func (s *Table[T]) One(ctx context.Context, filter filters.Filtered[T]) (T, error) {
	var zero T
	query, args := s.selectQuery(filter, 1)
	result, err := pgsql.ScanInto[T](s.queryable(ctx).QueryRow(ctx, query, args...), s.mapping)
	if err != nil {
		return zero, pgsql.TranslateError(err)
	}
	return result, nil
}

func (s *Table[T]) Many(ctx context.Context, filter filters.Filtered[T]) ([]T, error) {
	query, args := s.selectQuery(filter, 0)
	rows, err := s.queryable(ctx).Query(ctx, query, args...)
	if err != nil {
		return nil, pgsql.TranslateError(err)
	}
	defer rows.Close()

	results, err := pgsql.CollectRows[T](rows, s.mapping)
	if err != nil {
		return nil, pgsql.TranslateError(err)
	}
	return results, nil
}

func (s *Table[T]) DeleteOne(ctx context.Context, filter filters.Filtered[T]) error {
	where, args := pgsql.BuildWhere(filter, 0)
	query := fmt.Sprintf(
		"DELETE FROM %[1]s WHERE ctid IN (SELECT ctid FROM %[1]s WHERE %[2]s LIMIT 1)",
		s.name, where,
	)
	_, err := s.queryable(ctx).Exec(ctx, query, args...)
	return pgsql.TranslateError(err)
}

func (s *Table[T]) DeleteMany(ctx context.Context, filter filters.Filtered[T]) error {
	where, args := pgsql.BuildWhere(filter, 0)
	query := fmt.Sprintf("DELETE FROM %s WHERE %s", s.name, where)
	_, err := s.queryable(ctx).Exec(ctx, query, args...)
	return pgsql.TranslateError(err)
}
