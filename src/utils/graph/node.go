package graph

type Node[T any] struct {
	children []int
	data T
}

func NewNode[T any](data T) *Node[T] {
	return &Node[T]{
		data: data,
	}
}

func (n *Node[T]) AddChild(node int) {
	n.children = append(n.children, node)
}

func (n *Node[T]) Data() T {
	return n.data
}

func (n *Node[T]) Children() []int {
	return n.children
}