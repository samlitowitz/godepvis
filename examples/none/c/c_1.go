package c

import (
	"log"

	"github.com/samlitowitz/goimportcycle/examples/none/b"
)

func Fn1() {
	b.Fn()
	log.Println("C1")
}
