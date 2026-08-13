package postgre

import (
	"context"
	"fmt"
	"reflect"
	"strings"

	"github.com/Shitcode-Swamp/go-train-framework/src/core/persistance/postgre/internal/pgsql"
	"github.com/Shitcode-Swamp/go-train-framework/src/core/persistance/store"
)

func (s *Table[T]) Insert(ctx context.Context, doc T, options ...store.InsertOptions) (any, error) {
	return s.insert(ctx, s.pool, doc)
}

func (s *Table[T]) insert(ctx context.Context, q querier, doc T) (any, error) {
	cols, placeholders, args := s.insertColumns(doc)
	query := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)", s.name, strings.Join(cols, ", "), strings.Join(placeholders, ", "))

	if s.mapping.PKIndex < 0 {
		_, err := q.Exec(ctx, query, args...)
		return nil, pgsql.TranslateError(err)
	}
	return s.insertReturningID(ctx, q, query, args)
}

func (s *Table[T]) insertColumns(doc T) (cols, placeholders []string, args []any) {
	v := reflect.ValueOf(doc)
	cols = make([]string, len(s.mapping.Columns))
	placeholders = make([]string, len(s.mapping.Columns))
	args = make([]any, len(s.mapping.Columns))
	for i, c := range s.mapping.Columns {
		cols[i] = c.Column
		placeholders[i] = fmt.Sprintf("$%d", i+1)
		args[i] = v.Field(c.Index).Interface()
	}
	return cols, placeholders, args
}

func (s *Table[T]) insertReturningID(ctx context.Context, q querier, query string, args []any) (any, error) {
	query += fmt.Sprintf(" RETURNING %s", s.mapping.Columns[s.mapping.PKIndex].Column)
	var id any
	if err := q.QueryRow(ctx, query, args...).Scan(&id); err != nil {
		return nil, pgsql.TranslateError(err)
	}
	return id, nil
}

func (s *Table[T]) BulkInsert(ctx context.Context, docs []T, options ...store.InsertOptions) ([]any, error) {
	if len(docs) == 0 {
		return nil, nil
	}

	ids := make([]any, len(docs))
	err := s.withTx(ctx, func(q querier) error {
		for i, doc := range docs {
			id, err := s.insert(ctx, q, doc)
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
