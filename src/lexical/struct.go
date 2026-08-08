package lexical

import (
	"fmt"
	"strings"
)

type Field struct {
	Name string
	Type string
	Tag  string
}

func NewField(name, typ string) Field {
	return Field{Name: name, Type: typ}
}

func (f Field) WithTag(tag string) Field {
	f.Tag = tag
	return f
}

func (f Field) String() string {
	if f.Tag == "" {
		return fmt.Sprintf("%s %s", f.Name, f.Type)
	}
	return fmt.Sprintf("%s %s `%s`", f.Name, f.Type, f.Tag)
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
