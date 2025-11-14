package main

var VA = 1
var VB = 2
var (
	VC = 3
	VD = 4
)
var ve = 5

func VFn1() {
	var ve = 1
	_ = ve
}

func VFn2() {
	var ve = 1
	_ = ve
}

var _ = func() {
	var ve = 1
	_ = ve
}

var _ = func() {
	var ve = 1
	_ = ve
}
