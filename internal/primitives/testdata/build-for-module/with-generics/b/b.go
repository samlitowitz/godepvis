package b

import (
	"github.com/fake/fake/a"
)

func Fn() {}

var gtFn = a.IsGreater
var gtV = a.IsGreater(0, 1)
var st = a.Stack[int]{}
var sl = a.Slice[string, string]{}
var c = a.Clip([]int{0, 1, 2, 3})

func Sum[V a.Number](vs ...V) V {
	var s V
	for _, v := range vs {
		s += v
	}
	return s
}
