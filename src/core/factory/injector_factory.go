package factory

import (
	"go/ast"
	"strings"

	"github.com/Shitcode-Swamp/go-train-framework/src/core/injection"
	"github.com/Shitcode-Swamp/go-train-framework/src/utils/packages"
	utils_strings "github.com/Shitcode-Swamp/go-train-framework/src/utils/strings"
)

type Parser interface {
	ParseFuncParams(*ast.FuncType, *ast.File) []injection.Injectee
}

type Packages interface {
	GetPackage(filePath string) (string, bool)
}

type InjectorFactory struct {
	parser   Parser
	packages Packages
}

func NewInjectorFactory(parser Parser, pkgs Packages) *InjectorFactory {
	return &InjectorFactory{
		parser:   parser,
		packages: pkgs,
	}
}

func (f *InjectorFactory) Create(n ast.Node, file *ast.File, path string) (*injection.Injector, bool) {
	fDecl, ok := n.(*ast.FuncDecl)
	if !ok {
		return nil, false
	}
	tag := utils_strings.Extract(fDecl.Doc.Text(), `@Injection\(([^)]+)\)`)
	if tag == nil {
		return nil, false
	}
	token, flags := parseTag(*tag)
	pkg, ok := f.packages.GetPackage(path)
	if !ok {
		return nil, false
	}
	injectees := f.parser.ParseFuncParams(fDecl.Type, file)
	return injection.NewInjector(
		token,
		fDecl.Name.String(),
		functionReturnType(fDecl),
		fDecl.Name.String(),
		packages.Package{Name: file.Name.String(), Path: pkg},
		injectees,
		flags,
	), true
}

func parseTag(tag string) (string, []string) {
	parts := strings.Split(tag, ",")
	token := strings.TrimSpace(parts[0])
	if len(parts) == 1 {
		return token, nil
	}
	flags := make([]string, 0, len(parts)-1)
	for _, part := range parts[1:] {
		if flag := strings.TrimSpace(part); flag != "" {
			flags = append(flags, flag)
		}
	}
	return token, flags
}

func functionReturnType(fDecl *ast.FuncDecl) string {
	if fDecl.Type.Results == nil || len(fDecl.Type.Results.List) == 0 {
		return ""
	}
	result := fDecl.Type.Results.List[0]
	if result.Type == nil {
		return ""
	}
	return exprTypeName(result.Type)
}

func exprTypeName(expr ast.Expr) string {
	switch t := expr.(type) {
	case *ast.Ident:
		return t.Name
	case *ast.SelectorExpr:
		if x, ok := t.X.(*ast.Ident); ok {
			return x.Name + "." + t.Sel.Name
		}
	case *ast.StarExpr:
		return exprTypeName(t.X)
	}
	return ""
}
