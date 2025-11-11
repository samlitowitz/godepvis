package a

import (
	"log"

	"github.com/samlitowitz/godepvis/examples/interlinked/b"
)

func Fn() {
	log.Println("A")
	b.Fn()
}
