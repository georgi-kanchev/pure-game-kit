package example

import (
	"pure-game-kit/audio"
	"pure-game-kit/data/assets"
	"pure-game-kit/input/keyboard"
	"pure-game-kit/input/keyboard/key"
	"pure-game-kit/window"
)

func Audio() {
	var m1 = assets.LoadSound("examples/data/music2.ogg")
	var m2 = assets.LoadSound("examples/data/souls.ogg")
	var music = audio.New(m2)

	music.FadeIn = 1
	music.FadeOut = 1

	for window.KeepOpen() {
		if keyboard.IsKeyPressedOnce(key.A) {
			music.Play()
		}
		if keyboard.IsKeyPressedOnce(key.S) {
			music.LeftRight = 0
		}

		if music.IsFinishedOnce() {
			music.AssetId = m1
			music.Play()
		}
	}
}
