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
	window.Create("examples - audio", false, false)
	var sound1 = assets.LoadSound(3, "examples/data/door.ogg")
	var sound2 = assets.LoadSound(3, "examples/data/accept.mp3")
	var variation = assets.LoadAudioVariations(sound1, sound2)
	var audio = audio.New(variation)

	for window.KeepOpen() {
		if keyboard.IsKeyJustPressed(key.A) {
			audio.Play()
			fmt.Printf("play\n")
		}

		if keyboard.IsKeyJustPressed(key.S) {
			audio.Volume = 0.2
			audio.ApplyProperties()
			fmt.Printf("applied volume 20%%\n")
		}

		if audio.IsJustFinished() {
			fmt.Printf("just finished: %v\n", "true")
		}

		if keyboard.IsKeyJustPressed(key.Space) {
			if !audio.IsPaused() {
				audio.Pause()
				fmt.Printf("paused\n")
			} else {
				audio.Play()
				fmt.Printf("resumed\n")
			}
		}
	}
}
