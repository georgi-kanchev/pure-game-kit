package example

import (
	"fmt"
	"pure-kit/engine/execution/flow"
	"pure-kit/engine/graphics"
	"pure-kit/engine/utility/number"
	"pure-kit/engine/window"
)

func Flows() {
	var cam = graphics.NewCamera(1)
	var node = graphics.NewNode("", 0, 0)

	for window.KeepOpen() {
		cam.SetScreenAreaToWindow()
		cam.DrawNodes(&node)

		if flow.TrueEvery(1.0, "") {
			node.Angle = number.Wrap(node.Angle+45, 360)

			var lambda = flow.If(node.Angle == 45, "yes", "no")
			fmt.Printf("lambda angle is 45: %v (%v)\n", lambda, node.Angle)
		}
	}
}
