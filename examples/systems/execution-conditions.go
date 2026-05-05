package example

import (
	"pure-game-kit/packages/execution/condition"
	"pure-game-kit/packages/graphics"
	"pure-game-kit/packages/utility/number"
	"pure-game-kit/packages/utility/text"
	"pure-game-kit/packages/window"
)

func Conditions() {
	var cam = graphics.NewCamera(1)
	var quad = graphics.NewQuad(0, 0)

	for window.KeepOpen() {
		cam.DrawQuads(quad)

		if condition.TrueEvery(1.0, "") {
			quad.Angle = number.Wrap(quad.Angle+45, 0, 360)

			var lambda = condition.If(quad.Angle == 45, "yes", "no")
			print(text.New("lambda angle is 45: ", lambda, "(", quad.Angle, ")\n"))
		}
	}
}
