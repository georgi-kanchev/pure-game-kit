/*
A few helper functions, related to randomness. Also provides a way to use controlled randomness
in the form of seeds, as well as combining seeds. Has a few functions that act upon a collection -
shuffling it or choosing an item.
*/
package random

import (
	"fmt"
	"hash/fnv"
	"math"
	"math/rand/v2"
	"pure-game-kit/utility/number"
	"reflect"
)

func CombineSeeds[T number.Number](seeds ...T) T {
	if len(seeds) == 0 {
		var zero T
		switch any(zero).(type) {
		case float32:
			return T(number.NaN())
		case float64:
			return T(math.NaN())
		}
		return 0
	}

	var compact = make([]uint64, len(seeds))
	switch any(seeds[0]).(type) {
	case int, int8, int16, int32, int64:
		for i, s := range seeds {
			compact[i] = uint64(int64(s))
		}
	case uint, uint8, uint16, uint32, uint64:
		for i, s := range seeds {
			compact[i] = uint64(s)
		}
	case float32, float64:
		for i, s := range seeds {
			compact[i] = uint64(float64(s) * 1e9)
		}
	}

	var out = combineSeeds(compact...) // hash everything to one uint64
	var zero T                         // convert uint64 result back to T
	switch any(zero).(type) {
	case int:
		return T(int(out))
	case int8:
		return T(int8(out))
	case int16:
		return T(int16(out))
	case int32:
		return T(int32(out))
	case int64:
		return T(int64(out))
	case uint:
		return T(uint(out))
	case uint8:
		return T(uint8(out))
	case uint16:
		return T(uint16(out))
	case uint32:
		return T(uint32(out))
	case uint64:
		return T(uint64(out))
	case float32:
		return T(float32(out))
	case float64:
		return T(float64(out))
	}
	panic("unsupported type")
}

func Range[T number.Number](a, b T, seeds ...float32) T {
	switch any(a).(type) {
	case int, int8, int16, int32, int64:
		return T(rangeInt(int64(a), int64(b), seeds...))
	case uint, uint8, uint16, uint32, uint64:
		return T(rangeInt(uint64(a), uint64(b), seeds...))
	case float32, float64:
		return T(rangeFloat(float64(a), float64(b), seeds...))
	}
	panic("unsupported type")
}
func HasChance(percent float32, seeds ...float32) bool {
	if percent <= 0 {
		return false
	}
	return Range(float32(0), 100, seeds...) <= number.Smallest(100, percent)
}

func Shuffle[T any](items []T, seeds ...float32) []T {
	for i := len(items) - 1; i > 0; i-- {
		var j = Range(0, i, seeds...)
		items[i], items[j] = items[j], items[i]
	}
	return items
}
func Pick[T any](items ...T) T {
	return PickFrom(items)
}
func PickFrom[T any](items []T, seeds ...float32) T {
	return items[Range(0, len(items), seeds...)]
}

func Hash(value any) uint32 {
	var h = fnv.New32a()
	var val = reflect.ValueOf(value)
	if val.Kind() == reflect.Pointer {
		val = val.Elem()
	}
	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		if field.CanInterface() { // only exported fields
			h.Write(fmt.Appendf(nil, "%v", field.Interface()))
		}
	}
	return h.Sum32()
}

// =================================================================
// private

func hashSeed(seed, value uint64) uint64 {
	seed ^= value
	seed = (seed ^ (seed >> 16)) * 2246822519
	seed = (seed ^ (seed >> 13)) * 3266489917
	seed ^= seed >> 16
	return seed
}
func rangeInt[T number.Integer](a, b T, seeds ...float32) T {
	var ua, ub = uint64(a), uint64(b)
	if ua == ub {
		return a
	}
	if ua > ub {
		ua, ub = ub, ua
	}

	var diff = ub - ua
	var seed = CombineSeeds(seeds...)
	if seed != seed { // is NaN
		seed = rand.Float32()
	}
	var s = uint64(seed * 2147483647)
	s = (1103515245*s + 12345) % 2147483647
	var result = ua + (s*diff)/2147483647
	return T(result)
}
func rangeFloat[T number.Float](a, b T, seeds ...float32) T {
	var fa, fb = float64(a), float64(b)
	if fa == fb {
		return a
	}
	if fa > fb {
		fa, fb = fb, fa
	}

	var seed = CombineSeeds(seeds...)
	if seed != seed { // is NaN
		seed = rand.Float32()
	}

	var s = int(seed * 2147483647)
	s = (1103515245*s + 12345) % 2147483647
	var normalized = float64(s) / 2147483647.0
	var r = fa + (fb-fa)*normalized
	return T(r)
}
func combineSeeds(seeds ...uint64) uint64 {
	var seed = uint64(2654435769)
	for _, s := range seeds {
		seed = hashSeed(seed, s)
	}
	return seed
}
