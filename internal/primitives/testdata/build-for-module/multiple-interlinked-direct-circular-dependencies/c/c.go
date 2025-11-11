package c

import (
	"log"

	"github.com/fake/fake/b"
)

func Fn() {
	b.Fn()
	log.Println("C")
}
