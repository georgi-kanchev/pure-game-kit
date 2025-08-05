package example

import (
	"pure-kit/engine/data/assets"
	"pure-kit/engine/execution/flow"
	"pure-kit/engine/graphics"
	"pure-kit/engine/utility/color"
	"pure-kit/engine/window"
)

func Flows() {
	var cam = graphics.NewCamera(1)
	var node = graphics.NewNode(0, 0)
	var font = assets.LoadDefaultFont()
	var text = ""

	flow.NewSequence("flow-a",
		flow.WaitForDelay(2),
		flow.Do(func() {
			node.Angle = 45
			flow.Start("flow-b")
			text = "Flow A: step 1"
		}),
		flow.WaitForAnotherFlow("flow-b"),
		flow.WaitForDelay(2),
		flow.Do(func() {
			node.Angle = 0
			text = "Flow A: step 2"
		}),
		flow.WaitForDelay(2),
		flow.Do(func() {
			node.ScaleX, node.ScaleY = 3, 3
			text = "Flow A: step 3"
		}),
	)

	flow.NewSequence("flow-b",
		flow.WaitForDelay(2),
		flow.Do(func() {
			node.Color = color.Azure
			text = "Flow B: step 1"
		}),
		flow.WaitForDelay(2),
		flow.Do(func() {
			node.Color = color.Green
			text = "Flow B: step 2"
		}),
	)

	flow.Start("flow-a")

	for window.KeepOpen() {
		cam.SetScreenAreaToWindow()
		cam.PivotX, cam.PivotY = 0.5, 0.5
		cam.DrawNodes(&node)
		cam.PivotX, cam.PivotY = 0, 0
		cam.DrawText(font, text, 0, 0, 200, color.White)
	}
}
