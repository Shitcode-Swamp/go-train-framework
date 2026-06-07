package initialization

import (
	"errors"
	"fmt"

	"github.com/Shitcode-Swamp/go-train/src/core/injection"
	"github.com/Shitcode-Swamp/go-train/src/utils/graph"
)

var ErrorInjectorNotFound = errors.New("injector not found")

type Nodes []*graph.Node[*injection.Injector]

func (n Nodes) FindByToken(token string) (int, *graph.Node[*injection.Injector]) {
	for i, node := range n {
		nodeToken := node.Data().Token()
		if nodeToken != nil && *nodeToken == token {
			return i, node
		}
	}
	return -1, nil
}

type InitializationTree struct {
	nodes Nodes
}

type Visitor interface {
	Visit(injector *injection.Injector)
}

func NewInitializationTree() *InitializationTree {
	return &InitializationTree{
		nodes: make([]*graph.Node[*injection.Injector], 0),
	}
}

func (i *InitializationTree) Initialize(visitor Visitor) error {
	nodes, err := i.connectNodes()
	if err != nil {
		return err
	}
	for _, node := range nodes {
		visitor.Visit(node.Data())
	}
	return nil
}

func (i *InitializationTree) connectNodes() (Nodes, error) {
	for _, node := range i.nodes {
		err := i.addChildren(node)
		if err != nil {
			return nil, err
		}
	}
	return i.nodes, nil	
}

func (i *InitializationTree) addChildren(node *graph.Node[*injection.Injector]) error {
	for _, injectee := range node.Data().Injectees() {
		i, _ := i.nodes.FindByToken(injectee.Token())
		if i == -1 {
			return fmt.Errorf("%w: %s", ErrorInjectorNotFound, injectee.Token())
		}
		node.AddChild(i)
	}
	return nil
}

func (i *InitializationTree) Insert(injector *injection.Injector) {
	if injector == nil {
		return
	}
	i.nodes = append(i.nodes, graph.NewNode(injector))
}

