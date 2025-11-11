package b

import (
	"log"

	_ "github.com/fake/fake/a"
)

func Fn() {
	log.Println("B")
}
