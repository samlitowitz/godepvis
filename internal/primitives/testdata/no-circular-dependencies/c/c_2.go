package c

import (
	"log"

	"github.com/fake/no-circular-dependencies/b"
)

func Fn2() {
	b.Fn()
	log.Println("C2")
}
