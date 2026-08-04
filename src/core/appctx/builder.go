package appctx

import (
	injection_context "github.com/Shitcode-Swamp/go-train-framework/src/core/context"
	"github.com/Shitcode-Swamp/go-train-framework/src/core/factory"
	"github.com/Shitcode-Swamp/go-train-framework/src/core/initialization"
	"github.com/Shitcode-Swamp/go-train-framework/src/core/parsers"
	"github.com/Shitcode-Swamp/go-train-framework/src/core/render"
	"github.com/Shitcode-Swamp/go-train-framework/src/core/traverse"
	"github.com/Shitcode-Swamp/go-train-framework/src/core/writer"
)

const (
	contextType     = "AppContext"
	contextReceiver = "ctx"
)

type AppCtxBuilder struct {
	root   string
	output string
}

func NewAppCtxBuilder(root, output string) *AppCtxBuilder {
	return &AppCtxBuilder{
		root:   root,
		output: output,
	}
}

func (b *AppCtxBuilder) Build() error {
	tree := initialization.NewInitializationTree(nil)
	packages := parsers.NewPackages("./...")
	nodeFactory := factory.NewInjectorFactory(parsers.NewLexicalParser(), packages)
	fileTraverse := traverse.NewFileTraverse(tree, nodeFactory)
	if err := fileTraverse.Parse(b.root); err != nil {
		return err
	}
	if err := tree.Initialize(); err != nil {
		return err
	}
	renderer := render.NewInjectorRenderer(contextType, contextReceiver)
	visitor := injection_context.NewInjectionContextBuilder(renderer, contextType)
	if err := tree.Traverse(visitor); err != nil {
		return err
	}
	return writer.NewFileWriter(b.output).Write(visitor.Build())
}
