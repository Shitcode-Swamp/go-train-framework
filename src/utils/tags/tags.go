package tags

import (
	"go/ast"
	"reflect"
)

type Tag interface {
	Lookup(key string) (value string, ok bool)
}

type FieldsMeta map[string]Tag

func TagsFrom[T any]() FieldsMeta {
	var obj T
	reflected := reflect.TypeOf(&obj).Elem()
	tags := make(map[string]Tag, reflected.NumField())
	for i := range reflected.NumField() {
		field := reflected.Field(i)
		tags[field.Name] = field.Tag
	}
	return tags
}

func TagsFromAst(v *ast.StructType) FieldsMeta {
	fields := v.Fields.List
	tags := make(map[string]Tag, len(fields))
	for _, f := range fields {
		tags[f.Tag.Value] = &StringTag{ tag: f.Tag.Value }
	}
	return tags
}

type StringTag struct {
	tag string
}

func (t *StringTag) Lookup(key string) (value string, ok bool) {
	return "", false
}
