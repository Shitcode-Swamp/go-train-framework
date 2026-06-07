package main

import (
	"fmt"

	files_traverse "github.com/Shitcode-Swamp/go-train/src/core/files_traverse"
	"github.com/Shitcode-Swamp/go-train/src/core/injection"
	"github.com/Shitcode-Swamp/go-train/src/core/parsers"

	// "github.com/Shitcode-Swamp/go-train/src/core/injection"
	// "github.com/Shitcode-Swamp/go-train/src/core/someshi"
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
	tree := files_traverse.NewInitializationTree()
	traverse := files_traverse.NewFileTraverse(
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