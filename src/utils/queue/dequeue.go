package queue

type Dequeue[T any] []T 

func NewDequeue[T any](values ...T) Dequeue[T] {
	return Dequeue[T](values)
}

func (d *Dequeue[T]) Append(item T) {
	*d = append(*d, item)
}

func (d *Dequeue[T]) PopLeft() (T, bool) {
	if len(*d) == 0 {
		var zero T
		return zero, false
	}
	item := (*d)[0]
	*d = (*d)[1:]
	return item, true
}