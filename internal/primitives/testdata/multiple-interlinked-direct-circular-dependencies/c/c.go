package c

import (
	"log"

	"github.com/samlitowitz/godepvis/examples/interlinked/b"
)

func Fn() {
	b.Fn()
	log.Println("C")
}
