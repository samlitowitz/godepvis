package b

import (
	"github.com/fake/fake/a"
	"log"
)

func Fn() {
	a.Fn()
	log.Println("B")
}
