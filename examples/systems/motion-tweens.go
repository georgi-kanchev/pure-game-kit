package example

import (
	"pure-game-kit/packages/graphics"
	"pure-game-kit/packages/motion"
	"pure-game-kit/packages/utility/color/palette"
	"pure-game-kit/packages/utility/curve"
	"pure-game-kit/packages/utility/easing"
	"pure-game-kit/packages/window"
)

func Tweens() {
	var view = graphics.NewView(1)
	var angle = motion.NewTween(45).
		GoTo(2, easing.BounceOut, 360).
		GoTo(3, func(progress float32) float32 {
			var _, value = curve.Bezier(progress, 0, 0, 0.25, 1, 0.75, -0.5, 1, 1)
			return value
		}, 0)

	var position = motion.NewTween(-200, -200).
		GoTo(2, easing.ElasticOut, 200, 200).
		GoTo(3, easing.CubicOut, 0, 0).
		GoTo(2, easing.BackInOut, -200, 200)

	for window.KeepOpen() {
		view.MouseDragAndZoomSmoothly()

		view.DrawGrid(1, 100, 100, palette.Gray)

		var pos = position.CurrentValues()
		view.DrawQuad(pos[0], pos[1], 100, 100, angle.CurrentValues()[0], palette.White)

		if position.IsFinished() {
			position.Restart()
		}

		view.DrawTextDebug(true, true, true, true)
	}
}
