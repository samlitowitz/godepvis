package a

import (
	"log"

	"github.com/fake/fake/b"
)

func Fn() {
	log.Println("A")
	b.Fn()
}
