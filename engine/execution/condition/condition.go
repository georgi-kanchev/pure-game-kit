package condition

import "pure-kit/engine/internal"

func If[T any](condition bool, then, otherwise T) T {
	if condition {
		return then
	}
	return otherwise
}
func CallIf(condition bool, then func()) {
	if condition {
		then()
	}
}
func CallIfNotNil(function func()) {
	if function != nil {
		function()
	}
}
func CallIfElse(condition bool, then func(), otherwise func()) {
	if condition {
		then()
	}
	otherwise()
}
func CallAfter(seconds float32, function func()) {
	var t = internal.Runtime + seconds
	var _, has = internal.CallAfter[t]

	if !has {
		internal.CallAfter[t] = []func(){}
	}

	internal.CallAfter[t] = append(internal.CallAfter[t], function)
}
func CallFor(seconds float32, function func(remaining float32)) {
	var t = internal.Runtime + seconds
	var _, has = internal.CallFor[t]

	if !has {
		internal.CallFor[t] = []func(remaining float32){}
	}

	internal.CallFor[t] = append(internal.CallFor[t], function)
}

func TrueUponChange[T comparable](pointer *T) bool {
	var current = *pointer

	var prev, has = trueChanges[pointer]
	if !has || prev != current {
		trueChanges[pointer] = current
		return true
	}

	return false
}
func TrueOnce(condition bool, key any) bool {
	var prev = trueOnce[key]
	trueOnce[key] = condition
	return !prev && condition
}
func TrueEvery(seconds float32, key any) bool {
	var start, has = trueEvery[key]

	if !has || internal.Runtime > start+seconds {
		trueEvery[key] = internal.Runtime
		return has
	}

	return false
}

//=================================================================
// private

var trueChanges = make(map[any]any)
var trueOnce = make(map[any]bool)
var trueEvery = make(map[any]float32)
