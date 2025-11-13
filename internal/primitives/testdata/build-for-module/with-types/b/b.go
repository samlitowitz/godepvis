package b

import (
	"github.com/fake/fake/a"
	"log"
)

var aStruct = &a.AStruct{}
var aTStruct = &a.ATStruct{}
var aString = "Jello World"

func Fn() {
	log.Println(aString)
}
