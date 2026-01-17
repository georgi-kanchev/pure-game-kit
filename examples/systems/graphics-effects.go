package example

import (
	"pure-game-kit/data/assets"
	"pure-game-kit/debug"
	"pure-game-kit/graphics"
	"pure-game-kit/utility/number"
	"pure-game-kit/utility/time"
	"pure-game-kit/window"
)

func Effects() {
	var cam = graphics.NewCamera(4)
	var tex = assets.LoadTexture("examples/data/logo.PNG")
	var spr = graphics.NewSprite(tex, 0, 0)
	assets.SetTextureSmoothness(tex, true)

	spr.Effects = graphics.NewEffects()
	spr.Effects.Contrast = 0.7

	for window.KeepOpen() {
		cam.SetScreenAreaToWindow()
		cam.MouseDragAndZoomSmoothly()

		cam.DrawSprites(spr)

		spr.Effects.BlurX = 1 + number.Sine(time.Runtime())
		spr.Effects.BlurY = 1 + number.Cosine(time.Runtime())

		debug.Print(time.FrameRate())
	}
}
