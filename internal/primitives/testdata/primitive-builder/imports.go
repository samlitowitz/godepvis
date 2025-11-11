package main

import (
	"log"
	logalias "log"
	"os"
)

func main() {
	log.Println("log")
	logalias.Println("logalias")
	os.Exit(0)
}
