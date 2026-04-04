package example

import (
	"pure-game-kit/data/assets"
	"pure-game-kit/graphics"
	"pure-game-kit/utility/color/palette"
	"pure-game-kit/window"
)

func NinePatches() {
	var cam = graphics.NewCamera(1)
	var _, _, b = assets.LoadDefaultAtlasUI()
	var ninePatch = graphics.NewNinePatch(b[0], 0, 0)
	ninePatch.PivotX, ninePatch.PivotY = 0, 0
	ninePatch.Tint = palette.Cyan

	var bar = graphics.NewNinePatch(b[11], 0, 0)
	bar.PivotX, bar.PivotY = 0, 0

	for window.KeepOpen() {
		cam.DrawNinePatches(ninePatch, bar)

		var mx, my = ninePatch.PointToLocal(cam.MousePosition())
		ninePatch.Width, ninePatch.Height = mx, my
		bar.Width = mx
	}
}
