package b

import (
	"log"

	"github.com/samlitowitz/godepvis/v2/examples/simple/a"
)

func Fn() {
	a.Fn()
	log.Println("B")
}
