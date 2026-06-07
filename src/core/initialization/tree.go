package initialization

import (
	"errors"

	"github.com/Shitcode-Swamp/go-train/src/core/injection"
	"github.com/Shitcode-Swamp/go-train/src/utils/graph"
)

var ErrorInjectorNotFound = errors.New("injector not found")

type Nodes[T any] []T

type InitializationTree struct {
	tree *graph.Graph[*injection.Injector]
}

type Visitor interface {
	Visit(injector *injection.Injector)
}

func NewInitializationTree(nodes []*injection.Injector) *InitializationTree {
	if nodes == nil {
		nodes = make([]*injection.Injector, 0)
	}
	return &InitializationTree{
		tree: graph.NewGraph(nodes),
	}
}

func (it *InitializationTree) Initialize() error {
	return it.connectNodes()
}

func (it *InitializationTree) Traverse(visitor Visitor) error {
	nodes, err := it.tree.ReversedTopologicalSort()
	if err != nil {
		return err
	}
	for _, node := range nodes {
		visitor.Visit(node)
	}
	return nil
}

func (it *InitializationTree) connectNodes() error {
	tokenNodes := make(map[string]int)
	for idx, node := range it.tree.Nodes() {
		tokenNodes[*node.Token()] = idx
	}
	for _, node := range it.tree.Nodes() {
		if err := it.connectToChildren(node, tokenNodes); err != nil {
			return err
		}
	}
	return nil
}

func (it *InitializationTree) connectToChildren(
	node *injection.Injector,
	tokenNodes map[string]int,
) error {
	for _, injectee := range node.Injectees() {
		it.tree.AddEdge(tokenNodes[*node.Token()], tokenNodes[injectee.Token()])
	}
	return nil
}

func (it *InitializationTree) Insert(injector *injection.Injector) {
	if injector == nil {
		return
	}
	it.tree.AddNode(injector)
}
