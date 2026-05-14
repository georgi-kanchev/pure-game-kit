package main

import (
	"fmt"
	"pure-game-kit/packages/assets"
	"pure-game-kit/packages/audio"
	"pure-game-kit/packages/engine"
	"pure-game-kit/packages/input/keyboard"
	"pure-game-kit/packages/input/keyboard/key"
)

func main() {
	engine.Initialize("pure-game-kit", 60, 120, false, false)

	var sound = assets.LoadSound("examples/data/wood.mp3", 3)
	var audio = audio.New(sound)

	engine.Run(func() {
		if keyboard.IsKeyJustPressed(key.A) {
			audio.Pitch = 1.5
			audio.Play()
		}

		if keyboard.IsKeyJustPressed(key.S) {
			audio.Volume = 0.5
			audio.ApplyProperties()
		}

		if audio.IsJustFinished() {
			fmt.Printf("\"finished\": %v\n", "finished")
		}

		if keyboard.IsKeyJustPressed(key.Space) {
			if audio.IsPlaying() {
				audio.Pause()
			} else {
				audio.Resume()
			}
		}
	})
}
