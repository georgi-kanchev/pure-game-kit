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

// List wraps a standard slice to allow high-performance, allocation-free in-place mutations
// without cluttering the code with pointer syntax.
type List[T comparable] struct{ slice *[]T }

func NewList[T comparable](items ...T) List[T] {
	if len(items) == 0 {
		return List[T]{slice: &[]T{}}
	}
	var result = List[T]{slice: &[]T{}}
	result.Add(items...)
	return result
}
func NewListOfItem[T comparable](amount int, item T) List[T] {
	if amount <= 0 {
		return List[T]{slice: &[]T{}}
	}
	var result = List[T]{slice: &[]T{}}
	for range amount {
		result.Add(item)
	}
	return result
}
func NewListFromSlice[T comparable](slice *[]T) List[T] {
	return List[T]{slice: slice}
}

func (l List[T]) Clear()                       { *l.slice = (*l.slice)[:0] }
func (l List[T]) Add(items ...T)               { *l.slice = append(*l.slice, items...) }
func (l List[T]) Insert(index int, items ...T) { *l.slice = slices.Insert(*l.slice, index, items...) }
func (l List[T]) Reverse()                     { slices.Reverse(*l.slice) }
func (l List[T]) Set(index int, value T)       { (*l.slice)[index] = value }
func (l List[T]) Delete(index int)             { *l.slice = slices.Delete(*l.slice, index, index+1) }
func (l List[T]) Swap(indexA, indexB int) {
	(*l.slice)[indexA], (*l.slice)[indexB] = (*l.slice)[indexB], (*l.slice)[indexA]
}
func (l List[T]) Remove(items ...T) {
	for _, item := range items {
		for i, v := range *l.slice {
			if v == item {
				*l.slice = slices.Delete(*l.slice, i, i+1)
				break
			}
		}
	}
}
func (l List[T]) RemoveAt(indexes ...int) {
	// sort indexes descending so deleting early indices doesn't shift later ones we still need to delete
	slices.SortFunc(indexes, func(a, b int) int { return b - a })
	for _, index := range indexes {
		if index >= 0 && index < len(*l.slice) {
			*l.slice = slices.Delete(*l.slice, index, index+1)
		}
	}
}
func (l List[T]) RemoveUnordered(index int) {
	lastIdx := len(*l.slice) - 1 // fast but changes order
	(*l.slice)[index] = (*l.slice)[lastIdx]
	*l.slice = (*l.slice)[:lastIdx]
}
func (l List[T]) Shift(offset int) {
	var n = len(*l.slice)
	if n == 0 || offset == 0 {
		return
	}
	offset = ((offset % n) + n) % n // normalize offset

	var tmp = make([]T, n)
	copy(tmp[offset:], (*l.slice)[:n-offset])
	copy(tmp[:offset], (*l.slice)[n-offset:])
	copy(*l.slice, tmp)
}
func (l List[T]) ShiftIndexes(offset int, wrap bool, indexes ...int) {
	if len(*l.slice) == 0 || offset == 0 || len(indexes) == 0 {
		return
	}

	var n = len(*l.slice)
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
			tmp[i] = (*l.slice)[i]
			occupied[i] = true
		} else {
			tmp[target] = (*l.slice)[i]
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
		tmp[pos] = (*l.slice)[i]
		occupied[pos] = true
	}

	copy(*l.slice, tmp)
}
func (l List[T]) ShiftItems(offset int, wrap bool, items ...T) {
	var indexes = make([]int, 0, len(items))
	for _, item := range items {
		if idx := l.IndexOf(item); idx != -1 {
			indexes = append(indexes, idx)
		}
	}
	l.ShiftIndexes(offset, wrap, indexes...)
}
func (l List[T]) ShiftToFront(items ...T) {
	if len(items) == 0 || len(*l.slice) == 0 {
		return
	}

	var itemSet = make(map[T]struct{}, len(items)) // build a set for fast O(1) lookup
	for _, item := range items {
		itemSet[item] = struct{}{}
	}

	var writeIdx = 0 // filter out the moving items in-place to preserve remaining order
	for i := 0; i < len(*l.slice); i++ {
		if _, found := itemSet[(*l.slice)[i]]; !found {
			(*l.slice)[writeIdx] = (*l.slice)[i]
			writeIdx++
		}
	}
	*l.slice = (*l.slice)[:writeIdx]

	var origLen, itemsLen = len(*l.slice), len(items)
	*l.slice = slices.Grow(*l.slice, itemsLen) // ensure we have enough capacity
	*l.slice = (*l.slice)[:origLen+itemsLen]
	copy((*l.slice)[itemsLen:], (*l.slice)[:origLen]) // shift original items to the right
	copy(*l.slice, items)                             // copy shift items to the front
}
func (l List[T]) ShiftToEnd(items ...T) {
	if len(items) == 0 {
		return
	}

	for i := len(items) - 1; i >= 0; i-- {
		var block = items[i]
		if idx := l.IndexOf(block); idx != -1 { // find and remove the item if it exists
			*l.slice = slices.Delete(*l.slice, idx, idx+1)
		}
		*l.slice = append(*l.slice, block) // add to the end
	}
}
func (l List[T]) Join(lists ...List[T]) {
	var additionalLen = 0
	for _, arr := range lists {
		additionalLen += len(*arr.slice)
	}

	*l.slice = slices.Grow(*l.slice, additionalLen)
	for _, arr := range lists {
		*l.slice = append(*l.slice, *arr.slice...)
	}
}

func (l List[T]) Length() int           { return len(*l.slice) }
func (l List[T]) ToSlice() *[]T         { return (*[]T)(l.slice) }
func (l List[T]) First() T              { return (*l.slice)[0] }
func (l List[T]) Last() T               { return (*l.slice)[len(*l.slice)-1] }
func (l List[T]) Get(index int) T       { return (*l.slice)[index] }
func (l List[T]) IsEmpty() bool         { return len(*l.slice) == 0 }
func (l List[T]) Contains(value T) bool { return slices.Contains(*l.slice, value) }

func (l List[T]) Clone() { l.Copy() }
func (l List[T]) Copy() List[T] {
	var cloned = slices.Clone(*l.slice)
	return NewListFromSlice(&cloned)
}
func (l List[T]) IndexOf(value T) int {
	for i, v := range *l.slice {
		if v == value {
			return i
		}
	}
	return -1
}
func (l List[T]) HasDuplicates() bool {
	var seen = make(map[T]struct{}, len(*l.slice))
	for _, item := range *l.slice {
		if _, exists := seen[item]; exists {
			return true
		}
		seen[item] = struct{}{}
	}
	return false
}
func (l List[T]) String(divider string) string {
	builder.Reset()
	for i, elem := range *l.slice {
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
func MatrixString[T any](matrix [][]T, dividerRow, dividerColumn string) string {
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
