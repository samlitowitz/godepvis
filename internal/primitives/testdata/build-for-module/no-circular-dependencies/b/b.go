package b

import (
	"log"

	"github.com/fake/fake/a"
)

func Fn() {
	a.Fn()
	log.Println("B")
}
