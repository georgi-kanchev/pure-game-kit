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
	var nineSlice = graphics.NewNineSlice(t[32], 0, 0, [8]string{t[0], t[1], t[0], t[6], t[12], t[13], t[12], t[6]})
	nineSlice.Width, nineSlice.Height = 500, 500
	nineSlice.PivotX, nineSlice.PivotY = 0, 0
	nineSlice.SliceSizes = [4]float32{100, 100, 100, 100}
	nineSlice.SliceFlipX = [8]bool{false, false, true, true, true, false, false, false}

	for window.KeepOpen() {
		cam.SetScreenAreaToWindow()
		cam.DrawGrid(1, 100, color.Red)
		nineSlice.Color = color.Red
		cam.DrawNodes(&nineSlice.Node)
		nineSlice.Color = color.White
		cam.DrawNineSlices(&nineSlice)

		var mx, my = nineSlice.MousePosition(&cam)
		nineSlice.Width, nineSlice.Height = mx, my
	}
}
