package example

import (
	"fmt"
	"pure-game-kit/packages/assets"
	"pure-game-kit/packages/audio"
	"pure-game-kit/packages/input/keyboard"
	"pure-game-kit/packages/input/keyboard/key"
	"pure-game-kit/packages/window"
)

func Audio() {
	var sound = assets.LoadSound("examples/data/wood.mp3", 3)
	var audio = audio.New(sound)

	for window.KeepOpen() {
		if keyboard.IsKeyJustPressed(key.A) {
			audio.Play()
			fmt.Printf("play\n")
		}

		if keyboard.IsKeyJustPressed(key.S) {
			audio.Volume = 0.5
			audio.ApplyProperties()
			fmt.Printf("applied volume 50%%\n")
		}

		if audio.IsJustFinished() {
			fmt.Printf("just finished: %v\n", "true")
		}

		if keyboard.IsKeyJustPressed(key.Space) {
			if audio.IsPlaying() {
				audio.Pause()
				fmt.Printf("paused\n")
			} else {
				audio.Resume()
				fmt.Printf("resumed\n")
			}
		}
	}
}
