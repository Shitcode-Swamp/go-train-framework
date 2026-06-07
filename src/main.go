package main

import (
	"fmt"

	"github.com/Shitcode-Swamp/go-train/src/core/initialization"
	"github.com/Shitcode-Swamp/go-train/src/core/injection"
	"github.com/Shitcode-Swamp/go-train/src/core/parsers"
	"github.com/Shitcode-Swamp/go-train/src/core/traverse"

	"github.com/Shitcode-Swamp/go-train/src/utils/flags"
)

type Visitor struct {
}

func (v Visitor) Visit(injector *injection.Injector) {
	fmt.Printf("Visiting injector: %s\n", *injector.Token())
	injector.Expose()
}

func main() {
	flags := flags.NewFlags()
	flags.Parse()
	fmt.Println("root =", *flags.Root)
	tree := initialization.NewInitializationTree()
	traverse := traverse.NewFileTraverse(
		tree,
		parsers.NewLexicalParser(),
	)
	err := traverse.Parse(*flags.Root)
	if err != nil {
		fmt.Println("error: ", err)
	}
	err = tree.Initialize(Visitor{})
	if err != nil {
		fmt.Println("error: ", err)
	}
}
