package c

import (
	"log"

	"github.com/fake/fake/b"
	_ "github.com/fake/fake/b"
)

func Fn1() {
	b.Fn()
	log.Println("C1")
}
