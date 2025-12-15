package a

import (
	"log"

	"github.com/samlitowitz/godepvis/v3/examples/simple/b"
)

func Fn() {
	log.Println("A")
	b.Fn()
}
