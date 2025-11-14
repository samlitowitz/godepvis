package a

import (
	"log"

	"github.com/fake/fake/b"
)

type A struct{}

func (a *A) Fn(i any) {
	log.Println("A")
	bS, ok := i.(b.B)
	if !ok {
		return
	}
	bS.Fn(i)
}
