package example

import (
	"fmt"
	"pure-kit/engine/data/assets"
	"pure-kit/engine/graphics"
	"pure-kit/engine/scene"
	"pure-kit/engine/utility/color"
	"pure-kit/engine/window"
)

func Scenes() {
	var cam = graphics.NewCamera(7)
	var parent = graphics.NewSprite("map#1objects[1,0]", 0, 0)
	var data = assets.LoadTiledData("data/map.tmx")[0]
	var scene = scene.New(false, data)

	fmt.Printf("scene.BackgroundColor(): %v\n", scene.BackgroundColor())

	for window.KeepOpen() {
		cam.SetScreenAreaToWindow()
		cam.DrawGrid(1, 9, 9, color.Darken(color.Gray, 0.5))
		cam.DrawSprites(&parent)
	}
}
