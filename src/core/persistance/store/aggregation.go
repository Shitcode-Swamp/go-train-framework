package store

import "github.com/Shitcode-Swamp/go-train-framework/v3/src/core/persistance/filters"

type Stage interface {
	stage()
}

type FilterStage[T any] struct {
	Filter filters.Filtered[T]
}

func (FilterStage[T]) stage() {}

func Match[T any](f filters.Filtered[T]) Stage {
	return FilterStage[T]{Filter: f}
}

type ConvertStage struct {
	Field string
	To    string
}

func (ConvertStage) stage() {}

func Convert(field, to string) Stage {
	return ConvertStage{Field: field, To: to}
}

type GroupStage struct {
	By         []string
	Aggregates []AggregateExpr
}

func (GroupStage) stage() {}

func Group(by []string, aggregates ...AggregateExpr) Stage {
	return GroupStage{By: by, Aggregates: aggregates}
}

type AggregateExpr struct {
	Field string
	Op    AggregateOp
	As    string
}

type AggregateOp string

const (
	AggregateSum   AggregateOp = "sum"
	AggregateAvg   AggregateOp = "avg"
	AggregateMin   AggregateOp = "min"
	AggregateMax   AggregateOp = "max"
	AggregateCount AggregateOp = "count"
)

func Sum(field, as string) AggregateExpr {
	return AggregateExpr{Field: field, Op: AggregateSum, As: as}
}
func Avg(field, as string) AggregateExpr {
	return AggregateExpr{Field: field, Op: AggregateAvg, As: as}
}
func Min(field, as string) AggregateExpr {
	return AggregateExpr{Field: field, Op: AggregateMin, As: as}
}
func Max(field, as string) AggregateExpr {
	return AggregateExpr{Field: field, Op: AggregateMax, As: as}
}
func Count(as string) AggregateExpr { return AggregateExpr{Op: AggregateCount, As: as} }

type SortStage struct {
	By []SortField
}

func (SortStage) stage() {}

func Sort(by ...SortField) Stage {
	return SortStage{By: by}
}

type SortField struct {
	Field string
	Desc  bool
}

func Asc(field string) SortField  { return SortField{Field: field, Desc: false} }
func Desc(field string) SortField { return SortField{Field: field, Desc: true} }

type LimitStage struct {
	N int
}

func (LimitStage) stage() {}

func Limit(n int) Stage {
	return LimitStage{N: n}
}

type SkipStage struct {
	N int
}

func (SkipStage) stage() {}

func Skip(n int) Stage {
	return SkipStage{N: n}
}

type JoinStage struct {
	Collection   string
	LocalField   string
	ForeignField string
	As           string
}

func (JoinStage) stage() {}

func Join(collection, localField, foreignField, as string) Stage {
	return JoinStage{Collection: collection, LocalField: localField, ForeignField: foreignField, As: as}
}

type ProjectStage struct {
	Fields []ProjectField
}

func (ProjectStage) stage() {}

func Project(fields ...ProjectField) Stage {
	return ProjectStage{Fields: fields}
}

type ProjectField struct {
	From string
	To   string
}

func As(from, to string) ProjectField { return ProjectField{From: from, To: to} }
