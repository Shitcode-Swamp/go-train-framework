package filters

type Filtered[O any] struct {
	Op    Op
	Field string
	Value ValueExpr
	Subs  []Filtered[O]
}

func (f Filtered[O]) And(other Filtered[O]) Filtered[O] {
	return Filtered[O]{Op: OpAnd, Subs: []Filtered[O]{f, other}}
}

func (f Filtered[O]) Or(other Filtered[O]) Filtered[O] {
	return Filtered[O]{Op: OpOr, Subs: []Filtered[O]{f, other}}
}
