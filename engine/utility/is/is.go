package is

import (
	"fmt"
	"pure-kit/engine/utility/collection"
)

func AnyOf[T comparable](value T, values ...T) bool  { return collection.Contains(values, value) }
func OneOf[T comparable](value T, values ...T) bool  { return collection.Contains(values, value) }
func NoneOf[T comparable](value T, values ...T) bool { return !collection.Contains(values, value) }
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
