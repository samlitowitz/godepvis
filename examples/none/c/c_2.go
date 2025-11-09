package c

import (
	"log"

	"github.com/samlitowitz/goimportcycle/examples/none/b"
)

func Fn2() {
	b.Fn()
	log.Println("C2")
}
