package b

import (
	"log"

	"github.com/samlitowitz/godepvis/examples/transitive/a"
)

func Fn() {
	a.Fn()
	log.Println("B")
}
