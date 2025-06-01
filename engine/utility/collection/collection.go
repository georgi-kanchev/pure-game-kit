package collection

import (
	"fmt"
	"math"
	"pure-kit/engine/utility/number"
	"slices"
	"sort"
	"strings"
)

func Clone[T any](collection []T) []T {
	return slices.Clone(collection)
}
func SameItems[T any](amount int, item T) []T {
	result := make([]T, amount)
	for i := range amount {
		result[i] = item
	}
	return result
}

func Add[T any](collection *[]T, items ...T) {
	*collection = append(*collection, items...)
}
func Insert[T any](collection *[]T, index int, items ...T) {
	*collection = slices.Insert(*collection, index, items...)
}
func Remove[T comparable](collection *[]T, items ...T) {
	for _, item := range items {
		for i, v := range *collection {
			if v == item {
				*collection = slices.Delete(*collection, i, i+1)
				break
			}
		}
	}
}
func RemoveAt[T any](collection *[]T, indexes ...int) {
	// Sort indexes descending so deleting won't affect subsequent indices
	slices.SortFunc(indexes, func(a, b int) int { return a - b })
	for _, index := range indexes {
		if index >= 0 && index < len(*collection) {
			*collection = slices.Delete(*collection, index, index+1)
		}
	}
}

func IndexOf[T comparable](value T, collection []T) int {
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
	seen := make(map[T]struct{}, len(collection))
	for _, item := range collection {
		if _, exists := seen[item]; exists {
			return true
		}
		seen[item] = struct{}{}
	}
	return false
}

func Shift[T any](collection []T, offset int) {
	n := len(collection)
	if n == 0 || offset == 0 {
		return
	}
	offset = ((offset % n) + n) % n // normalize offset

	tmp := make([]T, n)
	copy(tmp[offset:], collection[:n-offset])
	copy(tmp[:offset], collection[n-offset:])
	copy(collection, tmp)
}
func ShiftIndexes[T any](collection []T, offset int, wrap bool, indexes ...int) {
	if len(collection) == 0 || offset == 0 || len(indexes) == 0 {
		return
	}

	n := len(collection)
	indexSet := make(map[int]bool, len(indexes))
	for _, idx := range indexes {
		if idx >= 0 && idx < n {
			indexSet[idx] = true
		}
	}

	// Sort indexes ascending for negative offset, descending for positive
	sorted := slices.Clone(indexes)
	if offset > 0 {
		sort.Sort(sort.Reverse(sort.IntSlice(sorted)))
	} else {
		sort.Ints(sorted)
	}

	tmp := make([]T, n)
	var zero T
	for i := range tmp {
		tmp[i] = zero
	}

	occupied := make(map[int]bool)

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

	pos := 0
	for i := 0; i < n; i++ {
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
	indexes := make([]int, 0, len(items))
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
func Surface[T comparable](collection []T, items []T) {
	if len(items) == 0 {
		return
	}

	for i := len(items) - 1; i >= 0; i-- {
		block := items[i]
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
func Sink[T comparable](collection []T, items []T) {
	if len(items) == 0 || len(collection) == 0 {
		return
	}

	// Build a set for faster lookup
	itemSet := make(map[T]struct{}, len(items))
	for _, item := range items {
		itemSet[item] = struct{}{}
	}

	// Step 1: Remove all items from collection
	dst := (collection)[:0]
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
	setA := make(map[T]struct{})
	for _, item := range collection {
		setA[item] = struct{}{}
	}

	result := make([]T, 0)
	seen := make(map[T]struct{})
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
	n := len(collection)
	if n == 0 {
		return nil
	}

	start = number.WrapInt(start, n)
	end = number.WrapInt(end, n)
	if start > end {
		start, end = end, start
	}

	result := make([]T, end-start)
	copy(result, collection[start:end])
	return result
}
func Join[T any](collection []T, otherCollections ...[]T) []T {
	totalLen := len(collection)
	for _, arr := range otherCollections {
		totalLen += len(arr)
	}

	result := make([]T, 0, totalLen)
	result = append(result, collection...)
	for _, arr := range otherCollections {
		result = append(result, arr...)
	}
	return result
}

func Rotate[T any](collection2D [][]T, direction int) [][]T {
	if direction == 0 || len(collection2D) == 0 || len(collection2D[0]) == 0 {
		return collection2D
	}

	dir := int(math.Abs(float64(direction))) % 4
	if dir == 0 {
		return collection2D
	}

	m, n := len(collection2D), len(collection2D[0])
	rotated := make([][]T, n)
	for i := range rotated {
		rotated[i] = make([]T, m)
	}

	if direction > 0 {
		for i := 0; i < n; i++ {
			for j := 0; j < m; j++ {
				rotated[i][j] = collection2D[m-j-1][i]
			}
		}
		return Rotate(rotated, direction-1)
	}

	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			rotated[i][j] = collection2D[j][n-i-1]
		}
	}
	return Rotate(rotated, direction+1)
}
func Flip[T any](collection2D [][]T, horizontally, vertically bool) [][]T {
	rows := len(collection2D)
	if rows == 0 {
		return collection2D
	}
	cols := len(collection2D[0])

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
	rows := len(collection2D)
	if rows == 0 {
		return nil
	}
	cols := len(collection2D[0])
	result := make([]T, 0, rows*cols)

	for i := 0; i < rows; i++ {
		result = append(result, collection2D[i]...)
	}
	return result
}

func ToText[T any](collection []T, divider string) string {
	var sb strings.Builder
	for i, elem := range collection {
		if i > 0 {
			sb.WriteString(divider)
		}
		sb.WriteString(fmt.Sprint(elem))
	}
	return sb.String()
}
func ToText2D[T any](collection2D [][]T, dividerRow, dividerColumn string) string {
	var sb strings.Builder
	for i, row := range collection2D {
		for j, elem := range row {
			sb.WriteString(fmt.Sprint(elem))
			if j < len(row)-1 {
				sb.WriteString(dividerRow)
			}
		}
		if i < len(collection2D)-1 {
			sb.WriteString(dividerColumn)
		}
	}
	return sb.String()
}
