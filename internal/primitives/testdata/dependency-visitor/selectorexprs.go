package main

import "fmt"

var sePrintln = fmt.Println

func sef1() {}

func sef2() {
	sef1()
	fmt.Println("sef2")
}

func init() {
	_, _ = sePrintln("init")
}
