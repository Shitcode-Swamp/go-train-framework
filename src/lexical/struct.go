package lexical

import (
	"fmt"
	"strings"
)

type Field struct {
	Name string
	Type string
}

func NewField(name, typ string) Field {
	return Field{Name: name, Type: typ}
}

func (f Field) String() string {
	return fmt.Sprintf("%s %s", f.Name, f.Type)
}

type Struct struct {
	Name   string
	Fields []Field
}

func NewStruct(name string, fields ...Field) Struct {
	return Struct{Name: name, Fields: fields}
}

func (s Struct) String() string {
	b := strings.Builder{}
	fmt.Fprintf(&b, "type %s struct {\n", s.Name)
	for _, field := range s.Fields {
		fmt.Fprintf(&b, "%s\n", field)
	}
	b.WriteString("}\n")
	return b.String()
}
