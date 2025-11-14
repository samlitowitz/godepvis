package b

import "github.com/fake/fake/a"

func Fn1() {
	var aS a.A
	_ = aS
	type b struct{}
	const c = 5
}

func Fn2() {
	var aS a.A
	_ = aS
	type b struct{}
	const c = 5
}

func Fn3() {
	func() {
		var aS a.A
		_ = aS
		type b struct{}
		const c = 5
	}()
}

var _ = func() {
	var aS a.A
	_ = aS
	type b struct{}
	const c = 5
}

var _ = func() {
	var aS a.A
	_ = aS
	type b struct{}
	const c = 5
}
