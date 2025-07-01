package main

import (
	"fmt"
	"pure-kit/engine/data/assets"
	"pure-kit/engine/motion/curves"
	"pure-kit/engine/motion/tween"
	"pure-kit/engine/render"
	"pure-kit/engine/utility/color"
	"pure-kit/engine/utility/seconds"
	"pure-kit/engine/window"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func main() {
	var cam = render.NewCamera()
	cam.Zoom = 3

	window.IsAntialiased = true
	var node = render.NewNode("tile", nil)
	node.AssetID = assets.LoadDefaultAtlasCursors()

	assets.LoadDefaultSoundsUserInterface()

	var chain = tween.From([]float32{0}).
		GoTo([]float32{1}, 1, curves.EaseLinear).
		CallWhenDone(func() {
			assets.PlaySound("ui-click")
		})

	for window.KeepOpen() {
		fmt.Printf("time.Runtime: %v\n", seconds.GetDelta())

		var w, h = window.Size()
		if rl.IsKeyPressed(rl.KeyA) {
			assets.PlaySound("ui-click")
			chain.Restart()
		}

		chain.Advance(float32(seconds.GetDelta()))

		cam.SetScreenArea(0, 0, w, h)
		cam.DrawGrid(1, 32, color.Darken(color.Gray, 0.5))
		cam.DrawNodes(&node)
	}
}
