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
	spr.Effects.Contrast = 0.7
	spr.Effects.SilhouetteColor = color.RGBA(255, 0, 0, 255)
	spr.Effects.PixelSize = 3
	spr.Effects.BlurX = 5
	spr.Effects.BlurY = 5

	for window.KeepOpen() {
		cam.SetScreenAreaToWindow()
		cam.MouseDragAndZoomSmoothly()

		cam.DrawSprites(spr)
	}
}
