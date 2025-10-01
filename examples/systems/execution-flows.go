package example

import (
	"fmt"
	"pure-kit/engine/data/assets"
	"pure-kit/engine/execution/condition"
	"pure-kit/engine/execution/flow"
	"pure-kit/engine/graphics"
	"pure-kit/engine/utility/color"
	"pure-kit/engine/utility/number"
	"pure-kit/engine/window"
)

func Flows() {
	var cam = graphics.NewCamera(1)
	var node = graphics.NewNode(0, 0)
	var font = assets.LoadDefaultFont()
	var text = ""
	var timerInt = 0

	flow.NewSequence("flow-a", true,
		flow.NowWaitForDelay(1),
		flow.NowDo(func() {
			node.Angle = 45
			flow.Run("flow-b")
			text = "Flow A: step 1"
		}),
		flow.NowWaitForFlow("flow-b"),
		flow.NowWaitForDelay(1),
		flow.NowDo(func() {
			node.Angle = 0
			text = "Flow A: step 2"
		}),
		flow.NowWaitForDelay(1),
		flow.NowDo(func() {
			node.ScaleX, node.ScaleY = 3, 3
			text = "Flow A: step 3"
		}),
		flow.NowDoLoop(5, func(i int) {
			fmt.Printf("loop i: %v\n", i)
		}),
		flow.NowDoLoop(number.ValueMaximum[int](), func(i int) {
			var timer = flow.CurrentStepTimer("flow-a")
			timerInt = int(timer)
			if timer == 0 {
				fmt.Printf("5 second timer started\n")
			}

			if timerInt > 0 && condition.TrueUponChange(&timerInt) {
				fmt.Printf("timer: %v\n", timer)
			}

			if timer > 5 {
				fmt.Printf("5 seconds timer ended\n")
				flow.GoToNextStep("flow-a")
			}
		}),
		flow.NowDo(func() {
			flow.Run("flow-a")
		}),
	)

	flow.NewSequence("flow-b", false,
		flow.NowWaitForDelay(1),
		flow.NowDo(func() {
			node.Color = color.Azure
			text = "Flow B: step 1"
		}),
		flow.NowWaitForDelay(1),
		flow.NowDo(func() {
			node.Color = color.Green
			text = "Flow B: step 2"
		}),
	)

	for window.KeepOpen() {
		cam.SetScreenAreaToWindow()
		cam.PivotX, cam.PivotY = 0.5, 0.5
		cam.DrawNodes(&node)
		cam.PivotX, cam.PivotY = 0, 0
		cam.DrawText(font, text, 0, 0, 200, color.White)
	}
}
