package main

import (
	"fmt"
	"strings"

	"github.com/Shitcode-Swamp/go-train/src/core/initialization"
	"github.com/Shitcode-Swamp/go-train/src/core/injection"
	"github.com/Shitcode-Swamp/go-train/src/core/parsers"
	"github.com/Shitcode-Swamp/go-train/src/core/traverse"

	"github.com/Shitcode-Swamp/go-train/src/utils/flags"
)

type Visitor struct {
	imports      []string
	declarations []string
}

func (v *Visitor) Visit(injector *injection.Injector) {
	v.imports = append(v.imports, injector.Import())
	v.declarations = append(v.declarations, injector.Init())
}

func (v Visitor) Print() {
	fmt.Println("Imports:")
	fmt.Println(strings.Join(v.imports, "\n"))
	fmt.Println("Declarations:")
	fmt.Println(strings.Join(v.declarations, "\n"))
}

func main() {
	flags := flags.NewFlags()
	flags.Parse()
	fmt.Println("root =", *flags.Root)
	tree := initialization.NewInitializationTree(nil)
	packages := parsers.NewPackages("./...")
	traverse := traverse.NewFileTraverse(
		tree,
		parsers.NewLexicalParser(),
		packages,
	)
	if err := traverse.Parse(*flags.Root); err != nil {
		fmt.Println("error: ", err)
	}
	visitor := &Visitor{ imports: make([]string, 0), declarations: make([]string, 0) }
	
	if err := tree.Initialize(); err != nil {
		fmt.Println("error: ", err)
	}
	tree.Traverse(visitor)
	visitor.Print()
}
