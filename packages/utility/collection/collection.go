// Helper functions for slices and maps. Some of them wrap standard functions to make them more
// digestible and clarify their API.
package collection

import (
	"cmp"
	"fmt"
	"pure-game-kit/packages/utility/number"
	"slices"
	"sort"
	"strings"
)

func Clear[T any](collection []T) []T {
	return collection[:0]
}
func Add[T any](collection []T, items ...T) []T {
	return append(collection, items...)
}
func Insert[T any](collection []T, index int, items ...T) []T {
	return slices.Insert(collection, index, items...)
}
func Reverse[T any](collection []T) {
	slices.Reverse(collection)
}
func Delete[T any](collection []T, index int) []T {
	return slices.Delete(collection, index, index+1)
}
func Swap[T any](collection []T, indexA, indexB int) {
	collection[indexA], collection[indexB] = collection[indexB], collection[indexA]
}
func Remove[T comparable](collection []T, items ...T) []T {
	for _, item := range items {
		for i, v := range collection {
			if v == item {
				collection = slices.Delete(collection, i, i+1)
				break
			}
		}
	}
	return collection
}
func RemoveAt[T any](collection []T, indexes ...int) []T {
	// sort indexes descending so deleting early indices doesn't shift later ones we still need to delete
	slices.SortFunc(indexes, func(a, b int) int { return b - a })
	for _, index := range indexes {
		if index >= 0 && index < len(collection) {
			collection = slices.Delete(collection, index, index+1)
		}
	}
	return collection
}
func RemoveUnordered[T any](collection []T, index int) []T {
	lastIdx := len(collection) - 1 // fast but changes order
	collection[index] = collection[lastIdx]
	return collection[:lastIdx]
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

	var sorted = slices.Clone(indexes)
	if offset > 0 { // sort indexes ascending for negative offset, descending for positive
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
		if idx := IndexOf(collection, item); idx != -1 {
			indexes = append(indexes, idx)
		}
	}
	ShiftIndexes(collection, offset, wrap, indexes...)
}
func ShiftToFront[T comparable](collection []T, items ...T) []T {
	if len(items) == 0 || len(collection) == 0 {
		return collection
	}

	var itemSet = make(map[T]struct{}, len(items)) // build a set for fast O(1) lookup
	for _, item := range items {
		itemSet[item] = struct{}{}
	}

	var writeIdx = 0 // filter out the moving items in-place to preserve remaining order
	for i := 0; i < len(collection); i++ {
		if _, found := itemSet[collection[i]]; !found {
			collection[writeIdx] = collection[i]
			writeIdx++
		}
	}
	collection = collection[:writeIdx]

	var origLen, itemsLen = len(collection), len(items)
	collection = slices.Grow(collection, itemsLen) // ensure we have enough capacity
	collection = collection[:origLen+itemsLen]
	copy(collection[itemsLen:], collection[:origLen]) // shift original items to the right
	copy(collection, items)                           // copy shift items to the front
	return collection
}
func ShiftToEnd[T comparable](collection []T, items ...T) []T {
	if len(items) == 0 {
		return collection
	}

	for i := len(items) - 1; i >= 0; i-- {
		var block = items[i]
		if idx := IndexOf(collection, block); idx != -1 { // find and remove the item if it exists
			collection = slices.Delete(collection, idx, idx+1)
		}
		collection = append(collection, block) // add to the end
	}
	return collection
}
func Join[T any](collection []T, collections ...[]T) []T {
	var additionalLen = 0
	for _, arr := range collections {
		additionalLen += len(arr)
	}

	collection = slices.Grow(collection, additionalLen)
	for _, arr := range collections {
		collection = append(collection, arr...)
	}
	return collection
}
func SameItems[T any](amount int, item T) []T {
	var result = make([]T, amount)
	for i := range amount {
		result[i] = item
	}
	return result
}

func Length[T any](collection []T) int {
	return len(collection)
}
func First[T any](collection []T) T {
	return collection[0]
}
func Last[T any](collection []T) T {
	return collection[len(collection)-1]
}
func IsEmpty[T any](collection []T) bool {
	return len(collection) == 0
}
func Contains[T comparable](collection []T, value T) bool {
	return slices.Contains(collection, value)
}

func Clone[T any](collection []T) []T {
	return Copy(collection)
}
func Copy[T any](collection []T) []T {
	return slices.Clone(collection)
}
func IndexOf[T comparable](collection []T, value T) int {
	for i, v := range collection {
		if v == value {
			return i
		}
	}
	return -1
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
func ToText[T any](collection []T, divider string) string {
	builder.Reset()
	for i, elem := range collection {
		if i > 0 {
			builder.WriteString(divider)
		}
		fmt.Fprint(&builder, elem)
	}
	return builder.String()
}

//=================================================================

func MatrixRotate[T any](matrix [][]T, direction int) [][]T {
	if direction == 0 || len(matrix) == 0 || len(matrix[0]) == 0 {
		return matrix
	}

	if number.Unsign(direction)%4 == 0 {
		return matrix
	}

	var m, n = len(matrix), len(matrix[0])
	var rotated = make([][]T, n)
	for i := range rotated {
		rotated[i] = make([]T, m)
	}

	if direction > 0 {
		for i := range n {
			for j := range m {
				rotated[i][j] = matrix[m-j-1][i]
			}
		}
		return MatrixRotate(rotated, direction-1)
	}

	for i := range n {
		for j := range m {
			rotated[i][j] = matrix[j][n-i-1]
		}
	}
	return MatrixRotate(rotated, direction+1)
}
func MatrixFlip[T any](matrix [][]T, horizontally, vertically bool) [][]T {
	var rows = len(matrix)
	if rows == 0 {
		return matrix
	}
	var cols = len(matrix[0])

	if horizontally {
		for i := range matrix {
			for j := 0; j < cols/2; j++ {
				matrix[i][j], matrix[i][cols-j-1] = matrix[i][cols-j-1], matrix[i][j]
			}
		}
	}

	if vertically {
		for i := range rows / 2 {
			matrix[i], matrix[rows-i-1] = matrix[rows-i-1], matrix[i]
		}
	}

	return matrix
}
func MatrixFlatten[T any](matrix [][]T) []T {
	var rows = len(matrix)
	if rows == 0 {
		return nil
	}
	var cols = len(matrix[0])
	var result = make([]T, 0, rows*cols)

	for i := range rows {
		result = append(result, matrix[i]...)
	}
	return result
}
func MatrixToText[T any](matrix [][]T, dividerRow, dividerColumn string) string {
	builder.Reset()
	for i, row := range matrix {
		for j, elem := range row {
			fmt.Fprint(&builder, elem)
			if j < len(row)-1 {
				builder.WriteString(dividerRow)
			}
		}
		if i < len(matrix)-1 {
			builder.WriteString(dividerColumn)
		}
	}
	return builder.String()
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
func MapClear[K comparable, V any](Map map[K]V) {
	clear(Map)
}

// private ========================================================

var builder strings.Builder
