package main

import (
	"fmt"
	"pure-kit/engine/geometry"
	"pure-kit/engine/utility/collection"
	"pure-kit/engine/window"
)

func main() {
	var items = [][]geometry.Point{
		{{X: 0, Y: 1}, {X: 3, Y: 4}, {X: 5, Y: 6}},
		{{X: 7, Y: 8}, {X: 9, Y: 10}, {X: 11, Y: 12}},
	}

	var grid = collection.ToText2D(items, ", ", "\n")

	fmt.Printf("%v", grid)

	for window.KeepOpen() {

	}
}
