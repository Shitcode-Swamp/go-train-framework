package contract

import "github.com/Shitcode-Swamp/go-train-framework/src/core/store/filters"

type Stage interface {
	stage()
}

type ConvertStage[T any] struct {
	Field filters.Field[T]
	To    string
}

func (ConvertStage[T]) stage() {}

type GroupStage[T any] struct {
	By         []filters.Field[T]
	Aggregates []AggregateExpr[T]
}

func (GroupStage[T]) stage() {}

type AggregateExpr[T any] struct {
	Field filters.Field[T]
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

type SortStage[T any] struct {
	By []SortField[T]
}

func (SortStage[T]) stage() {}

type SortField[T any] struct {
	Field filters.Field[T]
	Desc  bool
}

type LimitStage struct {
	N int
}

func (LimitStage) stage() {}

type SkipStage struct {
	N int
}

func (SkipStage) stage() {}

type JoinStage[T, F any] struct {
	Collection   string
	LocalField   filters.Field[T]
	ForeignField filters.Field[F]
	As           string
}

func (JoinStage[T, F]) stage() {}

type ProjectStage[T, U any] struct {
	Fields []ProjectField[T, U]
}

func (ProjectStage[T, U]) stage() {}

type ProjectField[T, U any] struct {
	From filters.Field[T]
	To   filters.Field[U]
}
