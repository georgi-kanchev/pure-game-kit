package example

import (
	"pure-kit/engine/audio"
	"pure-kit/engine/data/assets"
	"pure-kit/engine/input/keyboard"
	"pure-kit/engine/input/keyboard/key"
	"pure-kit/engine/window"
)

func Audio() {
	var sounds = assets.LoadDefaultSoundsUI()
	var m1 = assets.LoadMusic("examples/data/souls.ogg")
	// var m2 = assets.LoadSound("examples/data/music2.ogg")
	var music = audio.New(sounds[0])

	for window.KeepOpen() {
		if keyboard.IsKeyPressedOnce(key.A) {
			music.Play()
		}
		if keyboard.IsKeyPressedOnce(key.S) {
			music.AssetId = m1
		}
		if keyboard.IsKeyPressedOnce(key.D) {
			music.Volume = 0.5
		}
	}
}
