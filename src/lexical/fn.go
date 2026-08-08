package lexical

import (
	"fmt"
	"strings"
)

type Receiver struct {
	Var  string
	Type string
}

func NewReceiver(v, typ string) Receiver {
	return Receiver{Var: v, Type: typ}
}

func (r Receiver) String() string {
	return fmt.Sprintf("(%s %s)", r.Var, r.Type)
}

type Func struct {
	Receiver *Receiver
	Name     string
	Params   []Field
	Results  []string
	Body     []string
}

func NewFunc(name string, params []Field, results []string, body ...string) Func {
	return Func{Name: name, Params: params, Results: results, Body: body}
}

func (f Func) WithReceiver(receiver Receiver) Func {
	f.Receiver = &receiver
	return f
}

func (f Func) String() string {
	b := strings.Builder{}
	b.WriteString("func ")
	if f.Receiver != nil {
		fmt.Fprintf(&b, "%s ", f.Receiver)
	}
	fmt.Fprintf(&b, "%s(%s)", f.Name, f.paramsString())
	if results := f.resultsString(); results != "" {
		fmt.Fprintf(&b, " %s", results)
	}
	b.WriteString(" {\n")
	for _, line := range f.Body {
		fmt.Fprintf(&b, "%s\n", line)
	}
	b.WriteString("}\n")
	return b.String()
}

func (f Func) StringCall(args ...string) string {
	return fmt.Sprintf("%s(%s)", f.Name, strings.Join(args, ", "))
}

func (f Func) paramsString() string {
	params := make([]string, 0, len(f.Params))
	for _, param := range f.Params {
		params = append(params, param.String())
	}
	return strings.Join(params, ", ")
}

func (f Func) resultsString() string {
	switch len(f.Results) {
	case 0:
		return ""
	case 1:
		return f.Results[0]
	default:
		return fmt.Sprintf("(%s)", strings.Join(f.Results, ", "))
	}
}
