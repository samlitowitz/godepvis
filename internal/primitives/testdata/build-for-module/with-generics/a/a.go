package a

import "log"

func Fn() {
	log.Println("A")
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
