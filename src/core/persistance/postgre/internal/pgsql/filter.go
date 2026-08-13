package pgsql

import (
	"fmt"
	"strings"

	"github.com/Shitcode-Swamp/go-train-framework/src/core/persistance/filters"
)

type sqlBuilder struct {
	offset int
	args   []any
}

func (b *sqlBuilder) bind(value any) string {
	b.args = append(b.args, value)
	return fmt.Sprintf("$%d", b.offset+len(b.args))
}

func BuildWhere[T any](f filters.Filtered[T], offset int) (string, []any) {
	b := &sqlBuilder{offset: offset}
	clause := writeFiltered(f, b)
	if clause == "" {
		clause = "TRUE"
	}
	return clause, b.args
}

func writeFiltered[T any](f filters.Filtered[T], b *sqlBuilder) string {
	switch f.Op {
	case filters.OpAnd, filters.OpOr:
		joiner := " AND "
		if f.Op == filters.OpOr {
			joiner = " OR "
		}
		parts := make([]string, len(f.Subs))
		for i, sub := range f.Subs {
			parts[i] = "(" + writeFiltered(sub, b) + ")"
		}
		return strings.Join(parts, joiner)
	default:
		if f.Field == "" {
			return ""
		}
		return writeValueExpr(f.Field, f.Value, b)
	}
}

func writeValueExpr(field string, v filters.ValueExpr, b *sqlBuilder) string {
	switch v.Op {
	case filters.OpEq:
		return field + " = " + b.bind(unaryValue(v.Data))
	case filters.OpNe:
		return field + " <> " + b.bind(unaryValue(v.Data))
	case filters.OpGt:
		return field + " > " + b.bind(unaryValue(v.Data))
	case filters.OpGte:
		return field + " >= " + b.bind(unaryValue(v.Data))
	case filters.OpLt:
		return field + " < " + b.bind(unaryValue(v.Data))
	case filters.OpLte:
		return field + " <= " + b.bind(unaryValue(v.Data))
	case filters.OpIn:
		values := multiValues(v.Data)
		placeholders := make([]string, len(values))
		for i, val := range values {
			placeholders[i] = b.bind(val)
		}
		return field + " IN (" + strings.Join(placeholders, ", ") + ")"
	case filters.OpNot:
		return "NOT (" + writeValueExpr(field, v.Subs[0], b) + ")"
	default:
		return "TRUE"
	}
}

func unaryValue(data filters.Operable) any {
	return data.(filters.OperableUnary).Value
}

func multiValues(data filters.Operable) []any {
	return data.(filters.OperableMulti).Values
}

func SubstituteFields[T any](f filters.Filtered[T], resolve func(string) string) filters.Filtered[T] {
	if len(f.Subs) > 0 {
		subs := make([]filters.Filtered[T], len(f.Subs))
		for i, sub := range f.Subs {
			subs[i] = SubstituteFields(sub, resolve)
		}
		return filters.Filtered[T]{Op: f.Op, Subs: subs}
	}
	if f.Field == "" {
		return f
	}
	return filters.Filtered[T]{Op: f.Op, Field: resolve(f.Field), Value: f.Value}
}
