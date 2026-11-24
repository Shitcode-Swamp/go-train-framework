package storegen

import (
	"bytes"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"path"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

const storeTagKey = "store"

var (
	storeAnnotation = regexp.MustCompile(`@Store\(\s*"?([A-Za-z_][A-Za-z0-9_]*)"?\s*\)`)
	majorVersion    = regexp.MustCompile(`^v[0-9]+$`)
)

type entityParser struct {
	fset *token.FileSet
	file *ast.File
}

func parseEntityFile(filePath string) (*entityFile, error) {
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, filePath, nil, parser.ParseComments)
	if err != nil {
		return nil, err
	}
	p := &entityParser{fset: fset, file: file}
	return p.parse(filePath)
}

func (p *entityParser) parse(filePath string) (*entityFile, error) {
	result := newEntityFile(filePath, p.file.Name.Name)
	usedPackages := map[string]struct{}{}
	for _, decl := range p.file.Decls {
		gen, ok := decl.(*ast.GenDecl)
		if !ok || gen.Tok != token.TYPE {
			continue
		}
		for _, spec := range gen.Specs {
			entity, ok := p.parseEntity(gen, spec.(*ast.TypeSpec), usedPackages)
			if ok {
				result.entities = append(result.entities, entity)
			}
		}
	}
	result.imports = p.resolveImports(usedPackages)
	return result, nil
}

func (p *entityParser) parseEntity(gen *ast.GenDecl, spec *ast.TypeSpec, usedPackages map[string]struct{}) (storeEntity, bool) {
	name := storeEntityName(gen, spec)
	if name == "" {
		return storeEntity{}, false
	}
	structType, ok := spec.Type.(*ast.StructType)
	if !ok {
		return storeEntity{}, false
	}
	entity := storeEntity{name: name, typeName: spec.Name.Name}
	for _, field := range structType.Fields.List {
		entity.fields = append(entity.fields, p.parseFields(field, usedPackages)...)
	}
	return entity, true
}

func storeEntityName(gen *ast.GenDecl, spec *ast.TypeSpec) string {
	doc := spec.Doc
	if doc == nil && len(gen.Specs) == 1 {
		doc = gen.Doc
	}
	matches := storeAnnotation.FindStringSubmatch(doc.Text())
	if len(matches) < 2 {
		return ""
	}
	return matches[1]
}

func (p *entityParser) parseFields(field *ast.Field, usedPackages map[string]struct{}) []storeField {
	column := storeColumn(field)
	if column == "" || len(field.Names) == 0 {
		return nil
	}
	fieldType := p.typeString(field.Type)
	collectPackages(field.Type, usedPackages)
	fields := make([]storeField, 0, len(field.Names))
	for _, name := range field.Names {
		fields = append(fields, storeField{name: name.Name, typ: fieldType, column: column})
	}
	return fields
}

func storeColumn(field *ast.Field) string {
	if field.Tag == nil {
		return ""
	}
	raw, err := strconv.Unquote(field.Tag.Value)
	if err != nil {
		return ""
	}
	column, _, _ := strings.Cut(reflect.StructTag(raw).Get(storeTagKey), ",")
	if column == "-" {
		return ""
	}
	return column
}

func (p *entityParser) typeString(expr ast.Expr) string {
	var buf bytes.Buffer
	_ = printer.Fprint(&buf, p.fset, expr)
	return buf.String()
}

func collectPackages(expr ast.Expr, usedPackages map[string]struct{}) {
	ast.Inspect(expr, func(n ast.Node) bool {
		selector, ok := n.(*ast.SelectorExpr)
		if !ok {
			return true
		}
		if pkg, ok := selector.X.(*ast.Ident); ok {
			usedPackages[pkg.Name] = struct{}{}
		}
		return false
	})
}

func (p *entityParser) resolveImports(usedPackages map[string]struct{}) []entityImport {
	imports := make([]entityImport, 0, len(usedPackages))
	for _, spec := range p.file.Imports {
		importPath, err := strconv.Unquote(spec.Path.Value)
		if err != nil {
			continue
		}
		name := importName(spec, importPath)
		if _, ok := usedPackages[name]; ok {
			imports = append(imports, entityImport{alias: name, path: importPath})
		}
	}
	return imports
}

func importName(spec *ast.ImportSpec, importPath string) string {
	if spec.Name != nil {
		return spec.Name.Name
	}
	base := path.Base(importPath)
	if majorVersion.MatchString(base) {
		base = path.Base(path.Dir(importPath))
	}
	base = strings.TrimPrefix(base, "go-")
	base, _, _ = strings.Cut(base, ".")
	return strings.ReplaceAll(base, "-", "_")
}
