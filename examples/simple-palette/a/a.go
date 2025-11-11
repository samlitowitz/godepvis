package a

import (
	"log"

	"github.com/samlitowitz/godepvis/examples/simple/b"
)

func Fn() {
	log.Println("A")
	b.Fn()
}
