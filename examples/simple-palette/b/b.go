package b

import (
	"log"

	"github.com/samlitowitz/godepvis/v3/examples/simple/a"
)

func Fn() {
	a.Fn()
	log.Println("B")
}
