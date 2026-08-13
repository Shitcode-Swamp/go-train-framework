package aggregate

import (
	"fmt"

	"github.com/Shitcode-Swamp/go-train-framework/src/core/store/contract"
)

var aggregateFuncNames = map[contract.AggregateOp]string{
	contract.AggregateSum:   "SUM",
	contract.AggregateAvg:   "AVG",
	contract.AggregateMin:   "MIN",
	contract.AggregateMax:   "MAX",
	contract.AggregateCount: "COUNT",
}

func (b *aggregateQueryBuilder) applyConvert(st contract.ConvertStage) {
	b.fieldSources[st.Field] = b.resolve(st.Field) + "::" + st.To
}

func (b *aggregateQueryBuilder) applyGroup(st contract.GroupStage) {
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

func (b *aggregateQueryBuilder) aggregateExpr(agg contract.AggregateExpr) string {
	if agg.Op == contract.AggregateCount && agg.Field == "" {
		return "COUNT(*)"
	}
	return aggregateFuncNames[agg.Op] + "(" + b.resolve(agg.Field) + ")"
}

func (b *aggregateQueryBuilder) applySort(st contract.SortStage) {
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

func (b *aggregateQueryBuilder) applyLimit(st contract.LimitStage) {
	n := st.N
	b.limit = &n
}

func (b *aggregateQueryBuilder) applySkip(st contract.SkipStage) {
	n := st.N
	b.offset = &n
}

func (b *aggregateQueryBuilder) applyJoin(st contract.JoinStage) {
	on := fmt.Sprintf("%s.%s = %s.%s", b.table, st.LocalField, st.As, st.ForeignField)
	b.joins = append(b.joins, fmt.Sprintf("LEFT JOIN %s AS %s ON %s", st.Collection, st.As, on))
}

func (b *aggregateQueryBuilder) applyProject(st contract.ProjectStage) {
	for _, pf := range st.Fields {
		b.fieldSources[pf.To] = b.resolve(pf.From)
	}
}
