package c

import (
	"log"

	"github.com/samlitowitz/godepvis/examples/transitive/b"
)

func Fn() {
	b.Fn()
	log.Println("C")
}
