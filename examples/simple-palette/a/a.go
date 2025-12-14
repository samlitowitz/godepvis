package a

import (
	"log"

	"github.com/samlitowitz/godepvis/v2/examples/simple/b"
)

func Fn() {
	log.Println("A")
	b.Fn()
}
