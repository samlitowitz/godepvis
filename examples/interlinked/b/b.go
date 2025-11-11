package b

import (
	"log"

	"github.com/samlitowitz/godepvis/examples/interlinked/a"
	"github.com/samlitowitz/godepvis/examples/interlinked/c"
)

func Fn() {
	a.Fn()
	log.Println("B")
	c.Fn()
}
