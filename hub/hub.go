package main

import (
	"fmt"
	"pure-kit/engine/data/assets"
	"pure-kit/engine/graphics"
	"pure-kit/engine/utility/color"
	"pure-kit/engine/window"
)

func main() {
	var cam = graphics.NewCamera(7)
	var parent = graphics.NewNode("")
	var scenes = assets.LoadScenes("C:/Users/PC/Desktop/tiled/map.tmx")

	fmt.Printf("scenes: %v\n", scenes)

	for window.KeepOpen() {
		cam.SetScreenAreaToWindow()

		cam.DrawGrid(1, 9, color.Darken(color.Gray, 0.5))
		cam.DrawNodes(&parent)
	}
}
