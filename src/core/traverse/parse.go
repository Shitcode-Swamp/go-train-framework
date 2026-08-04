package traverse

import (
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"

	"github.com/Shitcode-Swamp/go-train/src/core/injection"
	utils_strings "github.com/Shitcode-Swamp/go-train/src/utils/strings"
)

type Injections = map[string]*injection.Injector

type Tree interface {
	Insert(injector *injection.Injector)
}

type Parser interface {
	ParseFuncParams(*ast.FuncType, *ast.File) []injection.Injectee
}

type Packages interface {
	GetPackage(filePath string) (string, bool)
}

type FileTraverse struct {
	collector Tree
	parser    Parser
	packages  Packages
}

func NewFileTraverse(
	collector Tree,
	parser Parser,
	packages Packages,
) *FileTraverse {
	return &FileTraverse{
		collector: collector,
		parser:    parser,
		packages:  packages,
	}
}

func (ft *FileTraverse) Parse(path string) error {
	root, err := filepath.Abs(path)
	if err != nil {
		return err
	}
	return ft.parseDir(root)
}

func (ft *FileTraverse) parseDir(path string) error {
	files, err := os.ReadDir(path)
	if err != nil {
		return err
	}
	for _, file := range files {
		err := ft.parseDirEntry(path, file)
		if err != nil {
			return err
		}
	}
	return nil
}

func (ft *FileTraverse) parseDirEntry(path string, entry os.DirEntry) error {
	if entry.IsDir() {
		return ft.parseDir(path + "/" + entry.Name())
	}
	if !strings.HasSuffix(entry.Name(), ".go") {
		return nil
	}
	return ft.parseFile(path + "/" + entry.Name())
}

func (ft *FileTraverse) parseFile(path string) error {
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, path, nil, parser.ParseComments)
	if err != nil {
		return err
	}
	ast.Inspect(file, func(n ast.Node) bool {
		injector := ft.inspectAstNode(n, file, path)
		if injector != nil {
			ft.collector.Insert(injector)
		}
		return true
	})
	return nil
}

func (ft *FileTraverse) inspectAstNode(n ast.Node, file *ast.File, path string) *injection.Injector {
	fDecl, ok := n.(*ast.FuncDecl)
	if !ok {
		return nil
	}
	token := utils_strings.Extract(fDecl.Doc.Text(), `@Injection\(([^)]+)\)`)
	if token == nil || fDecl.Type == nil {
		return nil
	}
	injectees := ft.parser.ParseFuncParams(fDecl.Type, file)
	pkg, ok := ft.packages.GetPackage(path)
	if !ok {
		return nil
	}
	return injection.NewInject(injectees).
		WithToken(*token).
		WithName(fDecl.Name.String()).
		WithTypeName(getFunctionReturnType(fDecl)).
		WithInitFn(fDecl.Name.String()).
		WithPackage(file.Name.String(), pkg)
}

func getFunctionReturnType(fDecl *ast.FuncDecl) string {
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
		return "" + exprTypeName(t.X)
	}
	return ""
}
