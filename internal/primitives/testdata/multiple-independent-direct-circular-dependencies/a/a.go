package a

import (
	"log"

	"github.com/samlitowitz/godepvis/examples/independent/b"
	"github.com/samlitowitz/godepvis/examples/independent/c"
)

func Fn() {
	log.Println("A")
	b.Fn()
	c.Fn()
}
