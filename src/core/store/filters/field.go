package filters

type Field[O, V any] string

func (f Field[O, V]) Is(vf Filter[V]) Filtered[O] {
	return Filtered[O]{Op: opLeaf, Field: string(f), Value: compile(vf)}
}
