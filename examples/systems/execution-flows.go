package example

import (
	"pure-game-kit/data/assets"
	"pure-game-kit/execution/flow"
	"pure-game-kit/graphics"
	"pure-game-kit/input/keyboard"
	"pure-game-kit/input/keyboard/key"
	"pure-game-kit/utility/color/palette"
	"pure-game-kit/utility/text"
	"pure-game-kit/window"
)

func Flows() {
	var cam = graphics.NewCamera(1)
	var node = graphics.NewNode(0, 0)
	var font = assets.LoadDefaultFont()
	var output = ""

	var second = flow.NewSequence()
	second.GoToStep(-1)
	second.SetSteps(
		flow.NowWaitForDelay(1),
		flow.NowDo(func() {
			node.Color = palette.Azure
			output = "Flow B: step 1"
		}),
		flow.NowWaitForDelay(1),
		flow.NowDo(func() {
			node.Color = palette.Green
			output = "Flow B: step 2"
		}))

	var first = flow.NewSequence()
	first.SetSteps(
		flow.NowWaitForDelay(1),
		flow.NowDo(func() {
			node.Angle = 45
			second.GoToStep(0)
			output = "Flow A: step 1"
		}),
		flow.NowWaitForSequence(second),
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
		flow.NowDoAndRepeat(50, func(i int) {
			output = text.New("looping: ", i)
		}),
		flow.NowDoAndRepeatFor(5, func() {
			output = text.New("timer: ", first.CurrentStepTimer())
		}),
		flow.NowDoAndKeepRepeating(func() {
			output = "Waiting for Space press"
			if keyboard.IsKeyJustPressed(key.Space) {
				first.GoToNextStep()
			}
		}),
		flow.NowDoAndRepeatFor(3, func() {
			output = "That's all! Restaring..."
		}),
		flow.NowDo(func() {
			first.GoToStep(0)
		}),
	)

	for window.KeepOpen() {
		cam.SetScreenAreaToWindow()
		cam.DrawNodes(&node)

		first.Update()
		second.Update()

		var x, y = cam.PointFromScreen(0, 0)
		cam.DrawText(font, output, x, y, 200, palette.White)

		if keyboard.IsKeyJustPressed(key.W) {
			first.Signal("W press")
		}
	}
}
