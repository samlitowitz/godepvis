package main

const CA = 1
const CB = 2
const (
	CC = 3
	CD = 4
)
const ce = 5

func CFn1() {
	const ce = 6
	var _ = func() {
		const ce = 6
	}
	func() {
		const ce = 6
	}()
}

func CFn2() {
	const ce = 6
}

var _ = func() {
	const ce = 6
}

var _ = func() {
	const ce = 6
}
