package main

import (
	"fmt"
	"pure-kit/engine/data/path"
	example "pure-kit/examples/systems"
)

func main() {
	fmt.Printf("path: %v\n", path.Folder(path.Executable()))

	example.Tiled()
}
