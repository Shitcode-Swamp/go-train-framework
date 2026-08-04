package main

import (
	"fmt"

	injection_context "github.com/Shitcode-Swamp/go-train/src/core/context"
	"github.com/Shitcode-Swamp/go-train/src/core/initialization"
	"github.com/Shitcode-Swamp/go-train/src/core/parsers"
	"github.com/Shitcode-Swamp/go-train/src/core/traverse"
	"github.com/Shitcode-Swamp/go-train/src/core/writer"

	"github.com/Shitcode-Swamp/go-train/src/utils/flags"
)

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
	visitor := injection_context.NewInjectionContextBuilder()
	if err := tree.Initialize(); err != nil {
		fmt.Println("error: ", err)
	}
	tree.Traverse(visitor)
	fileWriter := writer.NewFileWriter(*flags.Output)
	if err := fileWriter.Write(visitor.Build()); err != nil {
		fmt.Println("error: ", err)
	}
}
