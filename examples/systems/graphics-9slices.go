package example

import (
	"pure-kit/engine/data/assets"
	"pure-kit/engine/graphics"
	"pure-kit/engine/utility/color"
	"pure-kit/engine/window"
)

func NineSlices() {
	var cam = graphics.NewCamera(1)
	var _, t = assets.LoadDefaultAtlasUI()
	var asset = assets.LoadTextureNineSlice("button", [9]string{
		t[0], t[1], t[2],
		t[9], t[10], t[11],
		t[18], t[19], t[20]})
	var nineSlice = graphics.NewNineSlice(asset, 0, 0)
	nineSlice.PivotX, nineSlice.PivotY = 0, 0
	nineSlice.EdgeLeft = 100
	nineSlice.EdgeRight = 100
	nineSlice.EdgeBottom = 100
	nineSlice.EdgeTop = 100
	nineSlice.Color = color.Cyan

	for window.KeepOpen() {
		cam.SetScreenAreaToWindow()
		cam.DrawGrid(1, 100, 100, color.Red)
		cam.DrawNineSlices(&nineSlice)

		var mx, my = nineSlice.MousePosition(&cam)
		nineSlice.Width, nineSlice.Height = mx, my
	}
}
