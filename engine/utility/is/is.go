package is

import (
	"slices"

	"golang.org/x/exp/constraints"
)

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

func Any[T comparable](value T, values ...T) bool  { return slices.Contains(values, value) }
func One[T comparable](value T, values ...T) bool  { return slices.Contains(values, value) }
func None[T comparable](value T, values ...T) bool { return !slices.Contains(values, value) }
func All[T comparable](value T, values ...T) bool {
	for _, v := range values {
		if value != v {
			return false
		}
	}
	return true
}

func Flagged[T constraints.Integer](value, flag T) bool {
	return value&flag == flag
}
