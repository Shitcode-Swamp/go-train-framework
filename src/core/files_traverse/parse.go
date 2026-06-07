package files_traverse

import (
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"strings"

	"github.com/Shitcode-Swamp/go-train/src/core/injection"
	utils_strings "github.com/Shitcode-Swamp/go-train/src/utils/strings"
)

type Injections = map[string]*injection.Injector

type Collector interface {
	Insert(injector *injection.Injector)
}

type Parser interface {
	ParseFuncParams(*ast.FuncType, *ast.File) []injection.Injectee
}

type FileTraverse struct {
	collector Collector
	parser Parser
}

func NewFileTraverse(collector Collector, parser Parser) *FileTraverse {
	return &FileTraverse{
		collector: collector,
		parser: parser,
	}
}

func (ft *FileTraverse) Parse(path string) error {
	return ft.parseDir(path)
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
		return ft.parseDir(path+"/"+entry.Name())
	} 
	if !strings.HasSuffix(entry.Name(), ".go") {
		return nil
	}
	return ft.parseFile(path+"/"+entry.Name())
}

func (ft *FileTraverse) parseFile(path string) error {
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, path, nil, parser.ParseComments)
	if err != nil {
		return err
	}
	ast.Inspect(file, func(n ast.Node) bool {
		injector := ft.inspectAstNode(n, file)
		if injector != nil {
			ft.collector.Insert(injector)
		}
		return true
	})
	return nil
}

func (ft *FileTraverse) inspectAstNode(n ast.Node, file *ast.File) *injection.Injector {
	fDecl, ok := n.(*ast.FuncDecl)
	if !ok {
		return nil
	}
	token := utils_strings.Extract(fDecl.Doc.Text(), `@Injection\(([^)]+)\)`)
	if token == nil || fDecl.Type == nil {
		return nil
	}
	injectees := ft.parser.ParseFuncParams(fDecl.Type, file)
	return injection.NewInject(injectees).
		WithToken(*token).
		WithName(fDecl.Name.String()).
		WithPackage(file.Name.String())
}
