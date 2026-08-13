package pgsql

import (
	"reflect"

	"github.com/jackc/pgx/v5"
)

type Scanner interface {
	Scan(dest ...any) error
}

func ScanInto[T any](row Scanner, m TableMapping) (T, error) {
	var result T
	v := reflect.ValueOf(&result).Elem()
	dest := make([]any, len(m.Columns))
	for i, c := range m.Columns {
		dest[i] = v.Field(c.Index).Addr().Interface()
	}
	if err := row.Scan(dest...); err != nil {
		var zero T
		return zero, err
	}
	return result, nil
}

func CollectRows[T any](rows pgx.Rows, m TableMapping) ([]T, error) {
	var results []T
	for rows.Next() {
		item, err := ScanInto[T](rows, m)
		if err != nil {
			return nil, err
		}
		results = append(results, item)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return results, nil
}
