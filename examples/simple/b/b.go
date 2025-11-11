package b

import (
	"log"

	"github.com/samlitowitz/godepvis/examples/simple/a"
)

func Fn() {
	a.Fn()
	log.Println("B")
}
