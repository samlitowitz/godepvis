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

// subFnScopedSameDecl
// Test case for sub-function scoped variables with the same name
func subFnScopedSameDecl(i int) {
	if i > 0 {
		var bB any
		_ = bB
	}
	if i < 0 {
		var bB any
		_ = bB
	}
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
