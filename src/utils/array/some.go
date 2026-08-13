package array

func Some[T any](a []T, f func (e T) bool) bool {
	for _, e := range a {
		if f(e) {
			return true
		}
	}
	return false
}