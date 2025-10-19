package example

import (
	"fmt"
	"pure-game-kit/data/assets"
	"pure-game-kit/execution/condition"
	"pure-game-kit/execution/flow"
	"pure-game-kit/graphics"
	"pure-game-kit/input/keyboard"
	"pure-game-kit/input/keyboard/key"
	"pure-game-kit/utility/color"
	"pure-game-kit/utility/number"
	"pure-game-kit/window"
)

func Flows() {
	var cam = graphics.NewCamera(1)
	var node = graphics.NewNode(0, 0)
	var font = assets.LoadDefaultFont()
	var text = ""
	var timerInt = 0

	var b = flow.NewSequence()
	b.SetSteps(false,
		flow.NowWaitForDelay(1),
		flow.NowDo(func() {
			node.Color = color.Azure
			text = "Flow B: step 1"
		}),
		flow.NowWaitForDelay(1),
		flow.NowDo(func() {
			node.Color = color.Green
			text = "Flow B: step 2"
		}))

	var a = flow.NewSequence()
	a.SetSteps(true,
		flow.NowWaitForDelay(1),
		flow.NowDo(func() {
			node.Angle = 45
			b.Run()
			text = "Flow A: step 1"
		}),
		flow.NowWaitForSequence(b),
		flow.NowWaitForDelay(1),
		flow.NowDo(func() {
			node.Angle = 0
			text = "Flow A: step 2"
		}),
		flow.NowWaitForSignal("W press"),
		flow.NowWaitForDelay(1),
		flow.NowDo(func() {
			node.ScaleX, node.ScaleY = 3, 3
			text = "Flow A: step 3"
		}),
		flow.NowDoLoop(5, func(i int) {
			fmt.Printf("loop i: %v\n", i)
		}),
		flow.NowDoLoop(number.ValueMaximum[int](), func(i int) {
			var timer = a.CurrentStepTimer()
			timerInt = int(timer)
			if timer == 0 {
				fmt.Printf("5 second timer started\n")
			}

			if timerInt > 0 && condition.TrueUponChange(&timerInt) {
				fmt.Printf("timer: %v\n", timer)
			}

			if timer > 5 {
				fmt.Printf("5 seconds timer ended\n")
				a.GoToNextStep()
			}
		}),
		flow.NowDo(func() {
			a.Run()
		}),
	)

	for window.KeepOpen() {
		cam.SetScreenAreaToWindow()
		cam.PivotX, cam.PivotY = 0.5, 0.5
		cam.DrawNodes(&node)
		cam.PivotX, cam.PivotY = 0, 0
		cam.DrawText(font, text, 0, 0, 200, color.White)

		if keyboard.IsKeyPressedOnce(key.W) {
			a.Signal("W press")
		}
	}
}
