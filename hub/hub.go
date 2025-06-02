package main

import (
	"fmt"
	"pure-kit/engine/render"
	"pure-kit/engine/utility/color"
	"pure-kit/engine/window"
)

func main() {
	var camera = render.Camera{}

	fmt.Printf("camera: %v\n", camera)

	for window.KeepOpen() {
		camera.DrawRectangle(100, 100, 100, 100, color.Red)
	}
}
