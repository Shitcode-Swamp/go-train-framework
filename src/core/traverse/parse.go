package traverse

import (
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"

	"github.com/Shitcode-Swamp/go-train/src/core/injection"
)

type Tree interface {
	Insert(injector *injection.Injector)
}

type NodeFactory interface {
	Create(n ast.Node, file *ast.File, path string) (*injection.Injector, bool)
}

type FileTraverse struct {
	collector Tree
	factory   NodeFactory
}

func NewFileTraverse(
	collector Tree,
	factory NodeFactory,
) *FileTraverse {
	return &FileTraverse{
		collector: collector,
		factory:   factory,
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
		if injector, ok := ft.factory.Create(n, file, path); ok {
			ft.collector.Insert(injector)
		}
		return true
	})
	return nil
}
