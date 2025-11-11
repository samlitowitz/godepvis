package a

import (
	"log"

	"github.com/fake/fake/b"
	"github.com/fake/fake/c"
)

func Fn() {
	log.Println("A")
	b.Fn()
	c.Fn()
}
