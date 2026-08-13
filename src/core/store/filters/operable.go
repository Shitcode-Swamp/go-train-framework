package filters

type Operable interface {
	operable()
}

type OperableUnary struct {
	Value any
}

func (OperableUnary) operable() {}

type OperableMulti struct {
	Values []any
}

func (OperableMulti) operable() {}
