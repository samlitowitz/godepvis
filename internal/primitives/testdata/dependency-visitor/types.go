package main

type TA struct{}
type TB struct{}
type tc struct{}

func TFn1() {
	type tc struct{}
}

func TFn2() {
	type tc struct{}
}

var _ = func() {
	type tc struct{}
}

var _ = func() {
	type tc struct{}
}
