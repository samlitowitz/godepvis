package b

import (
	"log"

	"github.com/fake/fake/a"
)

type B struct{}

func (b *B) Fn() {
	aS := a.A{}
	aS.Fn()
	log.Println("B")
}
