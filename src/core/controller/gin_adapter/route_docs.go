package gin_adapter

import (
	"reflect"
	"strings"

	"github.com/Shitcode-Swamp/go-train-framework/src/core/controller/core"
)

type routeRegistry struct {
	routes []*routeDoc
}

func (reg *routeRegistry) add(method, path string) *routeDoc {
	rd := &routeDoc{method: method, path: path}
	reg.routes = append(reg.routes, rd)

	return rd
}

func swaggerPath(path string) string {
	segments := strings.Split(path, "/")

	for i, segment := range segments {
		if len(segment) > 1 && (segment[0] == ':' || segment[0] == '*') {
			segments[i] = "{" + segment[1:] + "}"
		}
	}

	return strings.Join(segments, "/")
}

type routeDoc struct {
	method string
	path   string

	query any
	body  any
	path_ any
}

func (r *routeDoc) RouteDocs() core.RouteDocs {
	return r
}

func (r *routeDoc) Query(query any) core.RouteDocs {
	r.query = query

	return r
}

func (r *routeDoc) Body(body any) core.RouteDocs {
	r.body = body

	return r
}

func (r *routeDoc) Path(path any) core.RouteDocs {
	r.path_ = path

	return r
}

func (r *routeDoc) operation() map[string]any {
	var params []map[string]any

	if r.path_ != nil {
		params = append(params, paramsFor(r.path_, "path")...)
	}

	if r.query != nil {
		params = append(params, paramsFor(r.query, "query")...)
	}

	if r.body != nil {
		params = append(params, map[string]any{
			"name":     "body",
			"in":       "body",
			"required": true,
			"schema":   schemaFor(reflect.TypeOf(r.body)),
		})
	}

	op := map[string]any{
		"responses": map[string]any{
			"200": map[string]any{"description": "OK"},
		},
	}

	if params != nil {
		op["parameters"] = params
	}

	return op
}

func paramsFor(v any, in string) []map[string]any {
	t := structTypeOf(v)
	if t == nil {
		return nil
	}

	params := make([]map[string]any, 0, t.NumField())

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		if !field.IsExported() {
			continue
		}

		params = append(params, map[string]any{
			"name":     jsonName(field),
			"in":       in,
			"required": in == "path",
			"type":     swaggerType(field.Type),
		})
	}

	return params
}

func schemaFor(t reflect.Type) map[string]any {
	for t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	switch t.Kind() {
	case reflect.Struct:
		properties := map[string]any{}

		for i := 0; i < t.NumField(); i++ {
			field := t.Field(i)
			if !field.IsExported() {
				continue
			}

			properties[jsonName(field)] = schemaFor(field.Type)
		}

		return map[string]any{"type": "object", "properties": properties}
	case reflect.Slice, reflect.Array:
		return map[string]any{"type": "array", "items": schemaFor(t.Elem())}
	default:
		return map[string]any{"type": swaggerType(t)}
	}
}

func structTypeOf(v any) reflect.Type {
	t := reflect.TypeOf(v)
	for t != nil && t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	if t == nil || t.Kind() != reflect.Struct {
		return nil
	}

	return t
}

func jsonName(field reflect.StructField) string {
	tag := field.Tag.Get("json")
	name, _, _ := strings.Cut(tag, ",")

	if name == "" || name == "-" {
		return field.Name
	}

	return name
}

func swaggerType(t reflect.Type) string {
	for t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	switch t.Kind() {
	case reflect.String:
		return "string"
	case reflect.Bool:
		return "boolean"
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return "integer"
	case reflect.Float32, reflect.Float64:
		return "number"
	default:
		return "string"
	}
}
