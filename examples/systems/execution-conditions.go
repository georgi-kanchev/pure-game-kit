package example

import (
	"pure-game-kit/execution/condition"
	"pure-game-kit/graphics"
	"pure-game-kit/utility/number"
	"pure-game-kit/utility/text"
	"pure-game-kit/window"
)

func Conditions() {
	var cam = graphics.NewCamera(1)
	var quad = graphics.NewQuad(0, 0)

	for window.KeepOpen() {
		cam.SetScreenAreaToWindow()
		cam.DrawQuads(quad)

		if condition.TrueEvery(1.0, "") {
			quad.Angle = number.Wrap(quad.Angle+45, 0, 360)

			var lambda = condition.If(quad.Angle == 45, "yes", "no")
			print(text.New("lambda angle is 45: ", lambda, "(", quad.Angle, ")\n"))
		}
	}
}
