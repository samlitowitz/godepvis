package b

import (
	"log"

	"github.com/fake/no-circular-dependencies/a"
)

func Fn() {
	a.Fn()
	log.Println("B")
}
