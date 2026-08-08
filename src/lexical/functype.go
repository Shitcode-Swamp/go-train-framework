package lexical

import (
	"fmt"
	"strings"
)

type FuncType struct {
	Name    string
	Params  []Field
	Results []string
}

func NewFuncType(name string, params []Field, results ...string) FuncType {
	return FuncType{Name: name, Params: params, Results: results}
}

func (t FuncType) String() string {
	b := strings.Builder{}
	fmt.Fprintf(&b, "type %s func(%s)", t.Name, fieldsString(t.Params))
	if results := resultsString(t.Results); results != "" {
		fmt.Fprintf(&b, " %s", results)
	}
	b.WriteString("\n")
	return b.String()
}
