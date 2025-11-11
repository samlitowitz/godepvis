package a

import (
	"log"

	"github.com/fake/fake/c"
)

func Fn() {
	log.Println("A")
	c.Fn()
}
