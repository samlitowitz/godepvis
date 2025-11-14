package a

import (
	"log"

	"github.com/fake/fake/b"
)

type A struct{}

func (a *A) Fn() {
	log.Println("A")
	bS := b.B{}
	bS.Fn()
}
