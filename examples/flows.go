package example

import (
	"pure-kit/engine/execution/flow"
	"pure-kit/engine/graphics"
	"pure-kit/engine/window"
)

func Flows() {
	var cam = graphics.NewCamera(1)
	var node = graphics.NewNode(0, 0)

	flow.NewSequence("my-new-flow",
		flow.Do(func() {
			println("first print")
			println("waiting for another flow...")
		}),
		flow.WaitForAnotherFlow("another-flow"),
		flow.Do(func() {
			println("resuming...")
			println("second print")
		}),
		flow.WaitForDelay(1),
		flow.Do(func() {
			println("third print")
		}),
	)

	flow.NewSequence("another-flow",
		flow.Do(func() {
			println("a whole another flow")
		}),
		flow.WaitForDelay(5),
		flow.Do(func() {
			println("another flow has finished")
		}),
	)

	flow.Start("my-new-flow")
	flow.Start("another-flow")

	for window.KeepOpen() {
		cam.SetScreenAreaToWindow()
		cam.DrawNodes(&node)
	}
}
