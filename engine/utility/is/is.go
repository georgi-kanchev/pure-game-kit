package is

import (
	"slices"

	"golang.org/x/exp/constraints"
)

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

func BitFlag[T constraints.Integer](value, flag T) bool {
	return value&flag == flag
}
