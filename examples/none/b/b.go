package b

import (
	"log"

	"github.com/samlitowitz/godepvis/examples/none/a"
)

func Fn() {
	a.Fn()
	log.Println("B")
}
