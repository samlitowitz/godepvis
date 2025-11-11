package b

import (
	"log"

	"github.com/fake/fake/a"
	"github.com/fake/fake/c"
)

func Fn() {
	a.Fn()
	log.Println("B")
	c.Fn()
}
