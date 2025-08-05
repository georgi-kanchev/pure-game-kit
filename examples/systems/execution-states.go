package example

import (
	"pure-kit/engine/execution/condition"
	"pure-kit/engine/execution/states"
	"pure-kit/engine/graphics"
	"pure-kit/engine/utility/seconds"
	"pure-kit/engine/window"
)

func States() {
	var cam = graphics.NewCamera(1)
	var node = graphics.NewNode(0, 0)

	states.NewMachine("my-state-machine",
		func() {
			node.X += seconds.FrameDelta() * 50
		},
		func() {
			node.Angle += seconds.FrameDelta() * 50
		},
		func() {
			node.ScaleX += seconds.FrameDelta()
			node.ScaleY += seconds.FrameDelta()
		},
	)

	condition.CallAfter(4, func() { states.GoToState("my-state-machine", 2) })
	condition.CallAfter(8, func() { states.GoToState("my-state-machine", 1) })

	for window.KeepOpen() {
		cam.SetScreenAreaToWindow()
		cam.DrawNodes(&node)
	}
}
