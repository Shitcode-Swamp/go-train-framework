package aggregate

import (
	"context"
	"fmt"
	"reflect"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/Shitcode-Swamp/go-train-framework/src/core/store/contract"
	"github.com/Shitcode-Swamp/go-train-framework/src/core/store/postgre/internal/pgsql"
)

type Pipeline[T, U any] struct {
	pool       *pgxpool.Pool
	table      string
	outMapping pgsql.TableMapping
}

func New[T, U any](pool *pgxpool.Pool, table string) contract.AggregateStore[T, U] {
	var zero U
	return &Pipeline[T, U]{
		pool:       pool,
		table:      table,
		outMapping: pgsql.MappingFor(reflect.TypeOf(zero)),
	}
}

var _ contract.AggregateStore[struct{}, struct{}] = (*Pipeline[struct{}, struct{}])(nil)

func (p *Pipeline[T, U]) Aggregate(ctx context.Context, pipeline []contract.Stage) ([]U, error) {
	qb := newAggregateQueryBuilder(p.table)
	for _, stage := range pipeline {
		switch st := stage.(type) {
		case contract.FilterStage[T]:
			applyFilterStage(qb, st)
		case contract.ConvertStage:
			qb.applyConvert(st)
		case contract.GroupStage:
			qb.applyGroup(st)
		case contract.SortStage:
			qb.applySort(st)
		case contract.LimitStage:
			qb.applyLimit(st)
		case contract.SkipStage:
			qb.applySkip(st)
		case contract.JoinStage:
			qb.applyJoin(st)
		case contract.ProjectStage:
			qb.applyProject(st)
		default:
			return nil, fmt.Errorf("%w: unrecognized aggregation stage %T", contract.ErrUnsupportedOperation, stage)
		}
	}

	query, args := qb.build(p.outMapping)
	rows, err := p.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, pgsql.TranslateError(err)
	}
	defer rows.Close()

	results, err := pgsql.CollectRows[U](rows, p.outMapping)
	if err != nil {
		return nil, pgsql.TranslateError(err)
	}
	return results, nil
}

func applyFilterStage[T any](b *aggregateQueryBuilder, fs contract.FilterStage[T]) {
	resolved := pgsql.SubstituteFields(fs.Filter, b.resolve)
	clause, args := pgsql.BuildWhere(resolved, len(b.args))
	b.args = append(b.args, args...)
	if b.grouped {
		b.having = append(b.having, clause)
	} else {
		b.where = append(b.where, clause)
	}
}
