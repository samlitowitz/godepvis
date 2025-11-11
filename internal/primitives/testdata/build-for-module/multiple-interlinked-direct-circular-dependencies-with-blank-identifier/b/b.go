package b

import (
	"log"

	_ "github.com/fake/fake/a"
	_ "github.com/fake/fake/c"
)

func Fn() {
	log.Println("B")
}
