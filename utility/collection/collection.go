/*
Helper functions for slices and maps. Some of them wrap standard functions to make them more
digestible and clarify their API.
*/
package collection

import (
	"cmp"
	"fmt"
	"pure-game-kit/utility/number"
	"pure-game-kit/utility/text"
	"slices"
	"sort"
)

func Clone[T any](collection []T) []T {
	return slices.Clone(collection)
}
func SameItems[T any](amount int, item T) []T {
	var result = make([]T, amount)
	for i := range amount {
		result[i] = item
	}
	return result
}

func Add[T any](collection []T, items ...T) []T {
	return append(collection, items...)
}
func Insert[T any](collection []T, index int, items ...T) []T {
	return slices.Insert(collection, index, items...)
}
func Remove[T comparable](collection []T, items ...T) []T {
	for _, item := range items {
		for i, v := range collection {
			if v == item {
				collection = slices.Delete(collection, i, i+1)
				break // Remove only the first match per item, like your pointer version
			}
		}
	}
	return collection
}
func RemoveAt[T any](collection []T, indexes ...int) []T {
	// Sort indexes descending so deletion doesn't affect subsequent indices
	var copy = Clone(collection)
	slices.SortFunc(indexes, func(a, b int) int { return b - a })
	for _, index := range indexes {
		if index >= 0 && index < len(copy) {
			copy = slices.Delete(copy, index, index+1)
		}
	}
	return copy
}

func IndexOf[T comparable](collection []T, value T) int {
	for i, v := range collection {
		if v == value {
			return i
		}
	}
	return -1
}
func IsEmpty[T any](collection T) bool {
	switch val := any(collection).(type) {
	case string:
		return val == ""
	case []any:
		return len(val) == 0
	case map[any]any:
		return len(val) == 0
	case chan any:
		return len(val) == 0
	default:
		return false
	}
}
func Contains[T comparable](collection []T, value T) bool {
	return slices.Contains(collection, value)
}
func HasDuplicates[T comparable](collection []T) bool {
	var seen = make(map[T]struct{}, len(collection))
	for _, item := range collection {
		if _, exists := seen[item]; exists {
			return true
		}
		seen[item] = struct{}{}
	}
	return false
}

func Shift[T any](collection []T, offset int) {
	var n = len(collection)
	if n == 0 || offset == 0 {
		return
	}
	offset = ((offset % n) + n) % n // normalize offset

	var tmp = make([]T, n)
	copy(tmp[offset:], collection[:n-offset])
	copy(tmp[:offset], collection[n-offset:])
	copy(collection, tmp)
}
func ShiftIndexes[T any](collection []T, offset int, wrap bool, indexes ...int) {
	if len(collection) == 0 || offset == 0 || len(indexes) == 0 {
		return
	}

	var n = len(collection)
	var indexSet = make(map[int]bool, len(indexes))
	for _, idx := range indexes {
		if idx >= 0 && idx < n {
			indexSet[idx] = true
		}
	}

	// Sort indexes ascending for negative offset, descending for positive
	var sorted = slices.Clone(indexes)
	if offset > 0 {
		sort.Sort(sort.Reverse(sort.IntSlice(sorted)))
	} else {
		sort.Ints(sorted)
	}

	var tmp = make([]T, n)
	var zero T
	for i := range tmp {
		tmp[i] = zero
	}

	var occupied = make(map[int]bool)

	for _, i := range sorted {
		target := i + offset
		if wrap {
			target = (target + n) % n
		} else {
			if target < 0 {
				target = 0
			} else if target >= n {
				target = n - 1
			}
		}

		if occupied[target] {
			tmp[i] = collection[i]
			occupied[i] = true
		} else {
			tmp[target] = collection[i]
			occupied[target] = true
		}
	}

	var pos = 0
	for i := range n {
		if _, moved := indexSet[i]; moved {
			continue
		}
		for occupied[pos] {
			pos++
		}
		tmp[pos] = collection[i]
		occupied[pos] = true
	}

	copy(collection, tmp)
}
func ShiftItems[T comparable](collection []T, offset int, wrap bool, items ...T) {
	var indexes = make([]int, 0, len(items))
	for _, item := range items {
		for i, val := range collection {
			if val == item {
				indexes = append(indexes, i)
				break
			}
		}
	}
	ShiftIndexes(collection, offset, wrap, indexes...)
}
func ShiftToEnd[T comparable](collection []T, items []T) {
	if len(items) == 0 {
		return
	}

	for i := len(items) - 1; i >= 0; i-- {
		var block = items[i]
		// Remove the item if it exists
		for j := range collection {
			if (collection)[j] == block {
				collection = slices.Delete((collection), j, j+1)
				break
			}
		}
		// Add to the end
		collection = append(collection, block)
	}
}
func ShiftToFront[T comparable](collection []T, items []T) {
	if len(items) == 0 || len(collection) == 0 {
		return
	}

	// Build a set for faster lookup
	var itemSet = make(map[T]struct{}, len(items))
	for _, item := range items {
		itemSet[item] = struct{}{}
	}

	// Step 1: Remove all items from collection
	var dst = (collection)[:0]
	for _, elem := range collection {
		if _, found := itemSet[elem]; !found {
			dst = append(dst, elem)
		}
	}
	collection = dst

	// Step 2: Prepend items in order (reversed, then re-reverse)
	for i := len(items) - 1; i >= 0; i-- {
		collection = append([]T{items[i]}, collection...)
	}
}

func Reverse[T any](collection []T) {
	for i, j := 0, len(collection)-1; i < j; i, j = i+1, j-1 {
		collection[i], collection[j] = collection[j], collection[i]
	}
}
func Overlap[T comparable](collection, otherCollection []T) []T {
	var setA = make(map[T]struct{})
	for _, item := range collection {
		setA[item] = struct{}{}
	}

	var result = make([]T, 0)
	var seen = make(map[T]struct{})
	for _, item := range otherCollection {
		if _, found := setA[item]; found {
			if _, already := seen[item]; !already {
				result = append(result, item)
				seen[item] = struct{}{}
			}
		}
	}
	return result
}
func Take[T any](collection []T, start, end int) []T {
	var n = len(collection)
	if n == 0 {
		return nil
	}

	start = number.Wrap(start, 0, n)
	end = number.Wrap(end, 0, n)
	if start > end {
		start, end = end, start
	}

	var result = make([]T, end-start)
	copy(result, collection[start:end])
	return result
}
func Join[T any](collections ...[]T) []T {
	var totalLen = len(collections)
	for _, arr := range collections {
		totalLen += len(arr)
	}

	var result = make([]T, 0, totalLen)
	for _, arr := range collections {
		result = append(result, arr...)
	}
	return result
}

func Rotate[T any](collection2D [][]T, direction int) [][]T {
	if direction == 0 || len(collection2D) == 0 || len(collection2D[0]) == 0 {
		return collection2D
	}

	if number.Unsign(direction)%4 == 0 {
		return collection2D
	}

	var m, n = len(collection2D), len(collection2D[0])
	var rotated = make([][]T, n)
	for i := range rotated {
		rotated[i] = make([]T, m)
	}

	if direction > 0 {
		for i := range n {
			for j := range m {
				rotated[i][j] = collection2D[m-j-1][i]
			}
		}
		return Rotate(rotated, direction-1)
	}

	for i := range n {
		for j := range m {
			rotated[i][j] = collection2D[j][n-i-1]
		}
	}
	return Rotate(rotated, direction+1)
}
func Flip[T any](collection2D [][]T, horizontally, vertically bool) [][]T {
	var rows = len(collection2D)
	if rows == 0 {
		return collection2D
	}
	var cols = len(collection2D[0])

	if horizontally {
		for i := range collection2D {
			for j := 0; j < cols/2; j++ {
				collection2D[i][j], collection2D[i][cols-j-1] = collection2D[i][cols-j-1], collection2D[i][j]
			}
		}
	}

	if vertically {
		for i := range rows / 2 {
			collection2D[i], collection2D[rows-i-1] = collection2D[rows-i-1], collection2D[i]
		}
	}

	return collection2D
}
func Flatten[T any](collection2D [][]T) []T {
	var rows = len(collection2D)
	if rows == 0 {
		return nil
	}
	var cols = len(collection2D[0])
	var result = make([]T, 0, rows*cols)

	for i := range rows {
		result = append(result, collection2D[i]...)
	}
	return result
}

func ToText[T any](collection []T, divider string) string {
	var builder = text.NewBuilder()
	for i, elem := range collection {
		if i > 0 {
			builder.WriteText(divider)
		}
		builder.WriteText(fmt.Sprint(elem))
	}
	return builder.ToText()
}
func ToText2D[T any](collection2D [][]T, dividerRow, dividerColumn string) string {
	var builder = text.NewBuilder()
	for i, row := range collection2D {
		for j, elem := range row {
			builder.WriteText(fmt.Sprint(elem))
			if j < len(row)-1 {
				builder.WriteText(dividerRow)
			}
		}
		if i < len(collection2D)-1 {
			builder.WriteText(dividerColumn)
		}
	}
	return builder.ToText()
}

func SortNumbers[T number.Number](collection ...T) {
	if len(collection) != 0 {
		slices.Sort(collection)
	}
}
func SortTexts(collection ...string) {
	if len(collection) != 0 {
		slices.Sort(collection)
	}
}
func SortByField[T any, F number.Number](s []T, field func(T) F) {
	slices.SortFunc(s, func(a, b T) int {
		return cmp.Compare(field(a), field(b))
	})
}

func MapKeys[K comparable, V any](Map map[K]V) []K {
	var keys = make([]K, 0, len(Map))
	for k := range Map {
		keys = append(keys, k)
	}
	return keys
}
func MapValues[K comparable, V any](Map map[K]V) []V {
	var values = make([]V, 0, len(Map))
	for _, v := range Map {
		values = append(values, v)
	}
	return values
}
