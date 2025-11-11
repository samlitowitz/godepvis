package c

import (
	"log"

	"github.com/fake/no-circular-dependencies/b"
)

func Fn1() {
	b.Fn()
	log.Println("C1")
}
