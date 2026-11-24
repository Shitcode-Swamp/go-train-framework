package postgre_aggregate

import (
	"fmt"

	"github.com/Shitcode-Swamp/go-train-framework/v3/src/core/persistance/store"
)

var aggregateFuncNames = map[store.AggregateOp]string{
	store.AggregateSum:   "SUM",
	store.AggregateAvg:   "AVG",
	store.AggregateMin:   "MIN",
	store.AggregateMax:   "MAX",
	store.AggregateCount: "COUNT",
}

func (b *aggregateQueryBuilder) applyConvert(st store.ConvertStage) {
	b.fieldSources[st.Field] = b.resolve(st.Field) + "::" + st.To
}

func (b *aggregateQueryBuilder) applyGroup(st store.GroupStage) {
	groupBy := make([]string, len(st.By))
	for i, field := range st.By {
		expr := b.resolve(field)
		groupBy[i] = expr
		b.fieldSources[field] = expr
	}
	for _, agg := range st.Aggregates {
		b.fieldSources[agg.As] = b.aggregateExpr(agg)
	}
	b.groupBy = groupBy
	b.grouped = true
}

func (b *aggregateQueryBuilder) aggregateExpr(agg store.AggregateExpr) string {
	if agg.Op == store.AggregateCount && agg.Field == "" {
		return "COUNT(*)"
	}
	return aggregateFuncNames[agg.Op] + "(" + b.resolve(agg.Field) + ")"
}

func (b *aggregateQueryBuilder) applySort(st store.SortStage) {
	parts := make([]string, len(st.By))
	for i, sf := range st.By {
		dir := "ASC"
		if sf.Desc {
			dir = "DESC"
		}
		parts[i] = b.resolve(sf.Field) + " " + dir
	}
	b.orderBy = parts
}

func (b *aggregateQueryBuilder) applyLimit(st store.LimitStage) {
	n := st.N
	b.limit = &n
}

func (b *aggregateQueryBuilder) applySkip(st store.SkipStage) {
	n := st.N
	b.offset = &n
}

func (b *aggregateQueryBuilder) applyJoin(st store.JoinStage) {
	on := fmt.Sprintf("%s.%s = %s.%s", b.table, st.LocalField, st.As, st.ForeignField)
	b.joins = append(b.joins, fmt.Sprintf("LEFT JOIN %s AS %s ON %s", st.Collection, st.As, on))
}

func (b *aggregateQueryBuilder) applyProject(st store.ProjectStage) {
	for _, pf := range st.Fields {
		b.fieldSources[pf.To] = b.resolve(pf.From)
	}
}
