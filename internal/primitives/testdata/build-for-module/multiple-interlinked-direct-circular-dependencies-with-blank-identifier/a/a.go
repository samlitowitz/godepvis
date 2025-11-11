package a

import (
	"log"

	_ "github.com/fake/fake/b"
)

func Fn() {
	log.Println("A")
}
