package example

import (
	"pure-game-kit/graphics"
	"pure-game-kit/motion"
	"pure-game-kit/motion/curve"
	"pure-game-kit/motion/easing"
	"pure-game-kit/utility/color/palette"
	"pure-game-kit/window"
)

func Tweens() {
	var cam = graphics.NewCamera(1)
	var angle = motion.NewTween(45).
		GoTo(2, easing.BounceOut, 360).
		GoTo(3, func(progress float32) float32 {
			var _, value = curve.Bezier(progress, [][2]float32{{0, 0}, {0.25, 1}, {0.75, -0.5}, {1, 1}})
			return value
		}, 0)

	var position = motion.NewTween(-200, -200).
		GoTo(2, easing.ElasticOut, 200, 200).
		GoTo(3, easing.CubicOut, 0, 0).
		GoTo(2, easing.BackInOut, -200, 200)

	for window.KeepOpen() {
		cam.SetScreenAreaToWindow()

		cam.DrawGrid(1, 100, 100, palette.Gray)

		var pos = position.CurrentValues()
		cam.DrawQuad(pos[0], pos[1], 100, 100, angle.CurrentValues()[0], palette.White)

		if position.IsFinished() {
			position.Restart()
		}
	}
}
