package filters

type Op string

const (
	opLeaf Op = "leaf"

	OpEq  Op = "eq"
	OpNe  Op = "ne"
	OpGt  Op = "gt"
	OpGte Op = "gte"
	OpLt  Op = "lt"
	OpLte Op = "lte"
	OpIn  Op = "in"
	OpNot Op = "not"

	OpAnd Op = "and"
	OpOr  Op = "or"
)

type Filter[V any] struct {
	Op    Op
	Data  Operable
	Inner *Filter[V]
}

func Eq[V any](value V) Filter[V]  { return Filter[V]{Op: OpEq, Data: OperableUnary{Value: value}} }
func Ne[V any](value V) Filter[V]  { return Filter[V]{Op: OpNe, Data: OperableUnary{Value: value}} }
func Gt[V any](value V) Filter[V]  { return Filter[V]{Op: OpGt, Data: OperableUnary{Value: value}} }
func Gte[V any](value V) Filter[V] { return Filter[V]{Op: OpGte, Data: OperableUnary{Value: value}} }
func Lt[V any](value V) Filter[V]  { return Filter[V]{Op: OpLt, Data: OperableUnary{Value: value}} }
func Lte[V any](value V) Filter[V] { return Filter[V]{Op: OpLte, Data: OperableUnary{Value: value}} }

func In[V any](values ...V) Filter[V] {
	anyValues := make([]any, len(values))
	for i, v := range values {
		anyValues[i] = v
	}
	return Filter[V]{Op: OpIn, Data: OperableMulti{Values: anyValues}}
}

func Not[V any](f Filter[V]) Filter[V] { return Filter[V]{Op: OpNot, Inner: &f} }
