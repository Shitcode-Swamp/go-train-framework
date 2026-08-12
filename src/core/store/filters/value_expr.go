package filters

type ValueExpr struct {
	Op   Op
	Data Operable
	Subs []ValueExpr
}

func compile[V any](f Filter[V]) ValueExpr {
	if f.Op == OpNot {
		return ValueExpr{Op: f.Op, Subs: []ValueExpr{compile(*f.Inner)}}
	}
	return ValueExpr{Op: f.Op, Data: f.Data}
}
