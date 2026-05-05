package example

import (
	"pure-game-kit/packages/assets"
	"pure-game-kit/packages/graphics"
	"pure-game-kit/packages/window"
)

func Effects() {
	var view = graphics.NewView(4)
	var tex = assets.LoadTexture("examples/data/logo.PNG")
	var spr = graphics.NewSprite(tex, 0, 0)
	assets.SetTextureSmoothness(tex, true)

	spr.Effects = graphics.NewEffects()
	spr.Effects.PixelSize = 3
	spr.Effects.DepthZ = 0.1
	spr.Effects.BlurX, spr.Effects.BlurY = 5, 5
	spr.ScaleX, spr.ScaleY = 0.2, 0.2

	var spr2 = graphics.NewSprite(tex, 50, 0)
	spr2.Effects = graphics.NewEffects()
	spr2.Effects.Saturation = 0.7
	spr2.Effects.DepthZ = 0.2
	spr2.ScaleX, spr2.ScaleY = 0.2, 0.2

	for window.KeepOpen() {
		view.MouseDragAndZoomSmoothly()

		view.DrawSprites(spr2, spr)
		view.DrawTextDebug(true, true, true, true)
	}
}
