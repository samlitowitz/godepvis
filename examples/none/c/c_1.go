package c

import (
	"log"

	"github.com/samlitowitz/godepvis/examples/none/b"
)

func Fn1() {
	b.Fn()
	log.Println("C1")
}
