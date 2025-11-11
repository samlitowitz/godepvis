package c

import (
	"log"

	"github.com/fake/no-circular-dependencies/b"
)

func Fn3() {
	b.Fn()
	log.Println("C3")
}
