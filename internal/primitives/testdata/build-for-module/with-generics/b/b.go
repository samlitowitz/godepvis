package b

import (
	"log"

	"github.com/fake/fake/a"
)

func Fn() {
	a.Fn()
	log.Println("B")
}

// SumNumbers sums the values of map m. It supports both integers
// and floats as map values.
func SumNumbers[K comparable, V a.Number](m map[K]V) V {
	var s V
	for _, v := range m {
		s += v

	}
	return s
}
