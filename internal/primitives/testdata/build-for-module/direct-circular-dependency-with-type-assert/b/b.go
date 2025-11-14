package b

import (
	"log"

	"github.com/fake/fake/a"
)

type B struct{}

func (b *B) Fn(i any) {
	aS, ok := i.(a.A)
	if !ok {
		return
	}
	aS.Fn(i)
	log.Println("B")
}
