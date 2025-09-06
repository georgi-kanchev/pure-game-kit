package is

import (
	"fmt"
	"slices"

	"golang.org/x/exp/constraints"
)

func AnyOf[T comparable](value T, values ...T) bool  { return slices.Contains(values, value) }
func OneOf[T comparable](value T, values ...T) bool  { return slices.Contains(values, value) }
func NoneOf[T comparable](value T, values ...T) bool { return !slices.Contains(values, value) }
func AllOf[T comparable](value T, values ...T) bool {
	for _, v := range values {
		if value != v {
			return false
		}
	}
	return true
}

func TypeOf(value any) string {
	return fmt.Sprintf("%T", value)
}

func BitFlag[T constraints.Integer](value, flag T) bool {
	return value&flag == flag
}
