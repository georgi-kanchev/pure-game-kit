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
	var node = graphics.NewNode(0, 0)

	for window.KeepOpen() {
		cam.SetScreenAreaToWindow()
		cam.DrawNodes(&node)

		if condition.TrueEvery(1.0, "") {
			node.Angle = number.Wrap(node.Angle+45, 0, 360)

			var lambda = condition.If(node.Angle == 45, "yes", "no")
			print(text.New("lambda angle is 45: ", lambda, "(", node.Angle, ")\n"))
		}
	}
}
