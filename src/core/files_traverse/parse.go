package files_traverse

import (
	"go/ast"
	"go/parser"
	"go/token"

	"github.com/Shitcode-Swamp/go-train/src/core/injection"
	"github.com/Shitcode-Swamp/go-train/src/core/parsers"
	utils_strings "github.com/Shitcode-Swamp/go-train/src/utils/strings"
)

type FileEntry struct {
	Injections map[string]*injection.Injector
}

func Parse(path string) (*FileEntry, error) {
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, path, nil, parser.ParseComments)
	if err != nil {
		return nil, err
	}
	injections := make(map[string]*injection.Injector)
	ast.Inspect(file, func(n ast.Node) bool {
		fDecl, ok := n.(*ast.FuncDecl)
		if !ok {
			return true
		}
		t := utils_strings.Extract(fDecl.Doc.Text(), `@Injection\(([^)]+)\)`)
		if t != nil && fDecl.Type != nil {
			injectees := parsers.ParseFuncParams(fDecl.Type, file)
			injections[*t] = injection.NewInject(injectees).
				WithName(fDecl.Name.String()).
				WithPackage(file.Name.String())
		}
		return true
	})

	return &FileEntry{Injections: injections}, nil
}
