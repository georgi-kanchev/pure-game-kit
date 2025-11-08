package example

import (
	"pure-game-kit/data/assets"
	"pure-game-kit/graphics"
	"pure-game-kit/utility/color"
	"pure-game-kit/window"
)

func Boxes() {
	var cam = graphics.NewCamera(1)
	var _, _, b = assets.LoadDefaultAtlasUI()
	var box = graphics.NewBox(b[0], 0, 0)
	box.PivotX, box.PivotY = 0, 0
	box.EdgeLeft = 100
	box.EdgeRight = 100
	box.EdgeBottom = 100
	box.EdgeTop = 100
	box.Color = color.Cyan

	var bar = graphics.NewBox(b[11], 0, 0)
	bar.PivotX, bar.PivotY = 0, 0
	bar.EdgeLeft = 100
	bar.EdgeRight = 100
	bar.EdgeBottom = 0
	bar.EdgeTop = 0

	for window.KeepOpen() {
		cam.SetScreenAreaToWindow()
		cam.DrawBoxes(&box, &bar)

		var mx, my = box.MousePosition(cam)
		box.Width, box.Height = mx, my
		bar.Width = mx
	}
}
