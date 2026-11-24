package graph

import (
	"errors"

	"github.com/Shitcode-Swamp/go-train-framework/v3/src/utils/queue"
)

var ErrGraphHasCycle = errors.New("Graph has cycle")

type Graph[T any] struct {
	nodes []T
	edges map[int][]int
}

func NewGraph[T any](nodes []T) *Graph[T] {
	return &Graph[T]{
		nodes: nodes,
		edges: make(map[int][]int),
	}
}

func (g *Graph[T]) AddNode(node T) {
	g.nodes = append(g.nodes, node)
}

func (g *Graph[T]) AddEdge(from, to int) {
	if from >= len(g.nodes) || to >= len(g.nodes) {
		return
	}
	g.edges[from] = append(g.edges[from], to)
}

func (g *Graph[T]) ReversedTopologicalSort() ([]T, error) {
	inDegree := g.inDegrees()
	independentVertecies := make([]int, 0)
	for i := range g.nodes {
		if inDegree[i] == 0 {
			independentVertecies = append(independentVertecies, i)
		}
	}
	zeroesQueue := queue.NewDequeue(independentVertecies...)
	sorted := make([]T, len(g.nodes))
	i := len(sorted) - 1
	for zero, ok := zeroesQueue.PopLeft(); ok; zero, ok = zeroesQueue.PopLeft() {
		sorted[i] = g.nodes[zero]
		i -= 1
		if len(sorted) > len(g.nodes) {
			return nil, ErrGraphHasCycle
		}
		for _, to := range g.edges[zero] {
			inDegree[to]--
			if inDegree[to] == 0 {
				zeroesQueue.Append(to)
			}
		}
	}
	return sorted, nil
}

func (g *Graph[T]) inDegrees() map[int]int {
	inDegree := make(map[int]int)
	for _, edges := range g.edges {
		for _, to := range edges {
			inDegree[to]++
		}
	}
	return inDegree
}

func (g *Graph[T]) Nodes() []T {
	return g.nodes
}

func (g *Graph[T]) Edges() map[int][]int {
	return g.edges
}
