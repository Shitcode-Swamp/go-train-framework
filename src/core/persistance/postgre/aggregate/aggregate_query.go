package aggregate

import (
	"fmt"
	"strings"

	"github.com/Shitcode-Swamp/go-train-framework/src/core/persistance/postgre/internal/pgsql"
)

type aggregateQueryBuilder struct {
	table        string
	joins        []string
	where        []string
	having       []string
	groupBy      []string
	orderBy      []string
	limit        *int
	offset       *int
	args         []any
	fieldSources map[string]string
	grouped      bool
}

func newAggregateQueryBuilder(table string) *aggregateQueryBuilder {
	return &aggregateQueryBuilder{table: table, fieldSources: map[string]string{}}
}

func (b *aggregateQueryBuilder) resolve(field string) string {
	if expr, ok := b.fieldSources[field]; ok {
		return expr
	}
	return field
}

func (b *aggregateQueryBuilder) build(outMapping pgsql.TableMapping) (string, []any) {
	selectList := make([]string, len(outMapping.Columns))
	for i, c := range outMapping.Columns {
		expr := b.resolve(c.Column)
		if expr == c.Column {
			selectList[i] = c.Column
		} else {
			selectList[i] = expr + " AS " + c.Column
		}
	}

	var sb strings.Builder
	sb.WriteString("SELECT ")
	sb.WriteString(strings.Join(selectList, ", "))
	sb.WriteString(" FROM ")
	sb.WriteString(b.table)
	for _, join := range b.joins {
		sb.WriteString(" ")
		sb.WriteString(join)
	}
	if len(b.where) > 0 {
		sb.WriteString(" WHERE ")
		sb.WriteString(joinParens(b.where, " AND "))
	}
	if len(b.groupBy) > 0 {
		sb.WriteString(" GROUP BY ")
		sb.WriteString(strings.Join(b.groupBy, ", "))
	}
	if len(b.having) > 0 {
		sb.WriteString(" HAVING ")
		sb.WriteString(joinParens(b.having, " AND "))
	}
	if len(b.orderBy) > 0 {
		sb.WriteString(" ORDER BY ")
		sb.WriteString(strings.Join(b.orderBy, ", "))
	}
	if b.limit != nil {
		sb.WriteString(fmt.Sprintf(" LIMIT %d", *b.limit))
	}
	if b.offset != nil {
		sb.WriteString(fmt.Sprintf(" OFFSET %d", *b.offset))
	}
	return sb.String(), b.args
}

func joinParens(clauses []string, sep string) string {
	wrapped := make([]string, len(clauses))
	for i, c := range clauses {
		wrapped[i] = "(" + c + ")"
	}
	return strings.Join(wrapped, sep)
}
