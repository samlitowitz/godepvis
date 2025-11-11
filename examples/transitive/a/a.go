package a

import (
	"log"

	"github.com/samlitowitz/godepvis/examples/transitive/c"
)

func Fn() {
	log.Println("A")
	c.Fn()
}
