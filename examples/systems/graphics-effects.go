package example

import (
	"pure-game-kit/data/assets"
	"pure-game-kit/graphics"
	"pure-game-kit/utility/color"
	"pure-game-kit/window"
)

func Effects() {
	var cam = graphics.NewCamera(4)
	var tex = assets.LoadTexture("examples/data/logo.PNG")
	var spr = graphics.NewSprite(tex, 0, 0)
	assets.SetTextureSmoothness(tex, true)

	spr.Effects = graphics.NewEffects()
	spr.Effects.SilhouetteColor = color.RGBA(255, 0, 0, 255)
	spr.Effects.PixelSize = 3
	spr.Effects.DepthZ = 0.1
	spr.Effects.BlurX, spr.Effects.BlurY = 5, 5
	spr.ScaleX, spr.ScaleY = 0.2, 0.2

	var spr2 = graphics.NewSprite("", 50, 0)
	spr2.Effects = graphics.NewEffects()
	spr2.Effects.Saturation = 0.7
	spr2.Effects.DepthZ = 0.2
	spr2.ScaleX, spr2.ScaleY = 0.2, 0.2
	spr2.Width, spr2.Height = 500, 500

	for window.KeepOpen() {
		cam.MouseDragAndZoomSmoothly()

		cam.DrawSprites(spr, spr2)
		cam.DrawTextDebug(true, true, true, true)
	}
}
