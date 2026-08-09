package filters

type Field[T any] string

type Filter[T any] interface {
	filter()
}

type EqFilter[T, V any] struct {
	Field Field[T]
	Value V
}

func (EqFilter[T, V]) filter() {}

func Eq[T, V any](field Field[T], value V) Filter[T] {
	return EqFilter[T, V]{Field: field, Value: value}
}

type NeFilter[T, V any] struct {
	Field Field[T]
	Value V
}

func (NeFilter[T, V]) filter() {}

func Ne[T, V any](field Field[T], value V) Filter[T] {
	return NeFilter[T, V]{Field: field, Value: value}
}

type GtFilter[T, V any] struct {
	Field Field[T]
	Value V
}

func (GtFilter[T, V]) filter() {}

func Gt[T, V any](field Field[T], value V) Filter[T] {
	return GtFilter[T, V]{Field: field, Value: value}
}

type GteFilter[T, V any] struct {
	Field Field[T]
	Value V
}

func (GteFilter[T, V]) filter() {}

func Gte[T, V any](field Field[T], value V) Filter[T] {
	return GteFilter[T, V]{Field: field, Value: value}
}

type LtFilter[T, V any] struct {
	Field Field[T]
	Value V
}

func (LtFilter[T, V]) filter() {}

func Lt[T, V any](field Field[T], value V) Filter[T] {
	return LtFilter[T, V]{Field: field, Value: value}
}

type LteFilter[T, V any] struct {
	Field Field[T]
	Value V
}

func (LteFilter[T, V]) filter() {}

func Lte[T, V any](field Field[T], value V) Filter[T] {
	return LteFilter[T, V]{Field: field, Value: value}
}

type InFilter[T, V any] struct {
	Field  Field[T]
	Values []V
}

func (InFilter[T, V]) filter() {}

func In[T, V any](field Field[T], values ...V) Filter[T] {
	return InFilter[T, V]{Field: field, Values: values}
}

type NotInFilter[T, V any] struct {
	Field  Field[T]
	Values []V
}

func (NotInFilter[T, V]) filter() {}

func NotIn[T, V any](field Field[T], values ...V) Filter[T] {
	return NotInFilter[T, V]{Field: field, Values: values}
}

type AndFilter[T any] struct {
	Filters []Filter[T]
}

func (AndFilter[T]) filter() {}

func And[T any](filters ...Filter[T]) Filter[T] {
	return AndFilter[T]{Filters: filters}
}

type OrFilter[T any] struct {
	Filters []Filter[T]
}

func (OrFilter[T]) filter() {}

func Or[T any](filters ...Filter[T]) Filter[T] {
	return OrFilter[T]{Filters: filters}
}

type NotFilter[T any] struct {
	Filter Filter[T]
}

func (NotFilter[T]) filter() {}

func Not[T any](f Filter[T]) Filter[T] {
	return NotFilter[T]{Filter: f}
}
