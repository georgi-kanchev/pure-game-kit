package main

import (
	"pure-kit/engine/data/assets"
	"pure-kit/engine/graphics"
	"pure-kit/engine/utility/collection"
	"pure-kit/engine/utility/color"
	"pure-kit/engine/window"
)

func main() {
	var cam = graphics.NewCamera(8)
	var textMap = map[[2]float32]string{
		{0, 0}: "#H", {1, 0}: "#e", {2, 0}: "#l", {3, 0}: "#l", {4, 0}: "#o", {4.7, 0}: "#,",
		{6, 0}: "#W", {7, 0}: "#o", {8, 0}: "#r", {9, 0}: "#l", {10, 0}: "#d", {11, 0}: "#!", {13, 0}: "#face-sad",
	}

	assets.LoadDefaultAtlasRetro()

	var textSymbols = collection.ToPointers(graphics.NewNodesGrid(textMap, 9, 9, nil))

	for window.KeepOpen() {
		cam.SetScreenAreaToWindow()

		cam.DrawGrid(1, 32, color.Darken(color.Gray, 0.5))
		cam.DrawNodes(textSymbols...)
	}
}
