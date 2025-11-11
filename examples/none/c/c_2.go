package c

import (
	"log"

	"github.com/samlitowitz/godepvis/examples/none/b"
)

func Fn2() {
	b.Fn()
	log.Println("C2")
}
