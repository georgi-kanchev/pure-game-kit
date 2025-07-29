package example

import (
	"pure-kit/engine/graphics"
	"pure-kit/engine/motion"
	"pure-kit/engine/motion/curves"
	"pure-kit/engine/utility/color"
	"pure-kit/engine/window"
)

func Tweens() {
	var angle = motion.NewTween(45).
		GoTo(2, curves.EaseBounceOut, 360).
		GoTo(3, func(progress float32) float32 {
			var _, value = curves.TraceBezier(progress, [][2]float32{{0, 0}, {0.25, 1}, {0.75, -0.5}, {1, 1}})
			return value
		}, 0)

	var position = motion.NewTween(-200, -200).
		GoTo(2, curves.EaseElasticOut, 200, 200).
		GoTo(3, curves.EaseCubicOut, 0, 0).
		GoTo(2, curves.EaseBackInOut, -200, 200)

	var cam = graphics.NewCamera(2)

	for window.KeepOpen() {
		cam.SetScreenAreaToWindow()

		var pos = position.CurrentValues()
		cam.DrawRectangle(pos[0], pos[1], 100, 100, angle.CurrentValues()[0], color.White)
	}
}
