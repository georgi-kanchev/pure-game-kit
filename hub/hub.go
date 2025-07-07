package main

import (
	"fmt"
	"pure-kit/engine/data/assets"
	"pure-kit/engine/graphics"
	"pure-kit/engine/scene"
	"pure-kit/engine/utility/color"
	"pure-kit/engine/window"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func main() {
	var cam = graphics.NewCamera(7)
	var parent = graphics.NewNode("map#1objects[1,0]")
	var data = assets.LoadTiledData("tiled/map.tmx")[0]
	var scene = scene.New(data)

	fmt.Printf("scene.BackgroundColor(): %v\n", scene.BackgroundColor())

	for window.KeepOpen() {
		cam.SetScreenAreaToWindow()

		if rl.IsKeyPressed(rl.KeyA) {
			scene.Unload()

		}

		cam.DrawGrid(1, 9, color.Darken(color.Gray, 0.5))
		cam.DrawNodes(&parent)
	}
}
