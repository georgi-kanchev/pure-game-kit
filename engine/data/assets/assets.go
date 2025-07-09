package assets

import (
	"pure-kit/engine/window"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func tryCreateWindow() {
	if !rl.IsWindowReady() {
		window.Recreate()
	}
}
func tryInitAudio() {
	if !rl.IsAudioDeviceReady() {
		rl.InitAudioDevice()
	}
}
