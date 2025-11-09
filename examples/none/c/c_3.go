package c

import (
	"log"

	"github.com/samlitowitz/goimportcycle/examples/none/b"
)

func Fn3() {
	b.Fn()
	log.Println("C3")
}
