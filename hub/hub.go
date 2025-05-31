package main

import (
	"fmt"
	animation "pure-tile-kit/engine/motion/animation"
	curves "pure-tile-kit/engine/motion/curves"
	tween "pure-tile-kit/engine/motion/tween"
	"pure-tile-kit/engine/window"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func main() {
	window.Color.R = 40

	var anim = animation.Sequence[string]{
		Items:          []string{"a", "b", "c", "d"},
		ItemsPerSecond: 1,
		IsLooping:      true,
		IsPlaying:      true,
	}

	var chain = tween.From([]float32{400, 200})
	chain.
		CallWhenDone(func() { println("start") }).
		GoTo([]float32{120, 120}, 3, curves.EaseBounceOut).
		CallWhenDone(func() { println("first") }).
		Wait(1).
		GoTo([]float32{400, 500}, 3, curves.EaseBackOut).
		CallWhileDoing(func(progress float32, current []float32) {
			fmt.Printf("progress: %v\n", progress)
		}).
		CallWhenDone(func() { println("second") }).
		Wait(1).
		GoTo([]float32{400, 900}, 3, curves.EaseElasticOut).
		CallWhenDone(func() { println("end") })

	for window.KeepOpen() {
		if rl.IsKeyPressed(rl.KeyA) {
			chain.Restart()
		}
		if rl.IsKeyPressed(rl.KeyS) {
			chain.Pause(true)
		}
		if rl.IsKeyPressed(rl.KeyD) {
			chain.Pause(false)
		}

		var frame, index = anim.Update(rl.GetFrameTime())
		var progress = float32(index) / float32(len(anim.Items))
		var result = fmt.Sprintf("%v: %v (%v %%)", index, *frame, progress*100)
		rl.DrawText(result, 0, 100, 64, rl.White)

		var current = chain.Update(rl.GetFrameTime())
		rl.DrawRectangle(int32(current[0]), int32(current[1]), 100, 100, rl.Red)
	}
}
