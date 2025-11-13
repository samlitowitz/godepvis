package a

func IsGreater[V float64 | int](a, b V) bool {
	return a > b
}

type Popper[T any] interface {
	Pop() T
}

type Stack[T any] struct {
	Data []T
}

func (m *Stack[T]) Push(item T) {
	m.Data = append(m.Data, item)
}

func (m *Stack[T]) Pop() T {
	item := m.Data[len(m.Data)-1]
	m.Data = m.Data[0 : len(m.Data)-1]
	return item
}

type Number interface {
	int64 | float64
}

type Slice[E, V any] []E

func (s Slice[E, V]) Map(iteratee func(E) V) []V {
	result := []V{}
	for _, item := range s {
		result = append(result, iteratee(item))
	}

	return result
}

func Clip[S ~[]E, E any](s S) S {
	return s[:len(s):len(s)]
}
