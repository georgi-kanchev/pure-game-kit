package example

import (
	"fmt"
	"pure-kit/engine/audio"
	"pure-kit/engine/data/assets"
	"pure-kit/engine/input/keyboard"
	"pure-kit/engine/input/keyboard/key"
	"pure-kit/engine/window"
)

func Audio() {
	var ui = assets.LoadDefaultSoundsUI()
	var asset = assets.LoadMusic("examples/data/music.ogg")
	var music = audio.New(asset)
	var woop = audio.New(ui[0])
	var hoop = audio.New(ui[4])

	for window.KeepOpen() {
		if keyboard.IsKeyPressedOnce(key.A) {
			music.Play()
		}
		if keyboard.IsKeyPressedOnce(key.S) {
			fmt.Printf("woop.IsPlaying(): %v\n", music.IsPlaying())
		}
		if keyboard.IsKeyPressedOnce(key.D) {
			woop.Play()
		}
		if keyboard.IsKeyPressedOnce(key.W) {
			hoop.Play()
		}
		fmt.Printf("hoop.IsPlaying(): %v\n", hoop.IsPlaying())
	}
}
