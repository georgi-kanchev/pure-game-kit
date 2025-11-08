package example

import (
	"pure-game-kit/data/assets"
	"pure-game-kit/execution/condition"
	"pure-game-kit/execution/flow"
	"pure-game-kit/graphics"
	"pure-game-kit/input/keyboard"
	"pure-game-kit/input/keyboard/key"
	"pure-game-kit/utility/color"
	"pure-game-kit/utility/number"
	"pure-game-kit/utility/text"
	"pure-game-kit/window"
)

func Flows() {
	var cam = graphics.NewCamera(1)
	var node = graphics.NewNode(0, 0)
	var font = assets.LoadDefaultFont()
	var output = ""
	var timerInt = 0

	var b = flow.NewSequence()
	b.SetSteps(false,
		flow.NowWaitForDelay(1),
		flow.NowDo(func() {
			node.Color = color.Azure
			output = "Flow B: step 1"
		}),
		flow.NowWaitForDelay(1),
		flow.NowDo(func() {
			node.Color = color.Green
			output = "Flow B: step 2"
		}))

	var a = flow.NewSequence()
	a.SetSteps(true,
		flow.NowWaitForDelay(1),
		flow.NowDo(func() {
			node.Angle = 45
			b.Run()
			output = "Flow A: step 1"
		}),
		flow.NowWaitForSequence(b),
		flow.NowWaitForDelay(1),
		flow.NowDo(func() {
			node.Angle = 0
			output = "Flow A: step 2 (waiting for W press)"
		}),
		flow.NowWaitForSignal("W press"),
		flow.NowWaitForDelay(0.1),
		flow.NowDo(func() {
			node.ScaleX, node.ScaleY = 3, 3
			output = "Flow A: step 3"
		}),
		flow.NowDoLoop(50, func(i int) {
			output = text.New("looping: ", i)
		}),
		flow.NowDoLoop(number.ValueMaximum[int](), func(i int) {
			var timer = a.CurrentStepTimer()
			timerInt = int(timer)
			if timer == 0 {
				output = text.New("5 second timer started")
			}

			if timerInt > 0 && condition.JustChanged(&timerInt) {
				output = text.New("timer: ", timer)
			}

			if timer > 5 {
				output = text.New("5 seconds timer ended")
				a.GoToNextStep()
			}
		}),
		flow.NowWaitForDelay(1),
		flow.NowDo(func() {
			a.Run()
		}),
	)

	for window.KeepOpen() {
		cam.SetScreenAreaToWindow()
		cam.PivotX, cam.PivotY = 0.5, 0.5
		cam.DrawNodes(&node)
		cam.PivotX, cam.PivotY = 0, 0
		cam.DrawText(font, output, 0, 0, 200, color.White)

		if keyboard.IsKeyJustPressed(key.W) {
			a.Signal("W press")
		}
	}
}
