package assets

import (
	"pure-game-kit/data/path"
	"pure-game-kit/internal"
	"pure-game-kit/utility/text"
	"pure-game-kit/window"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func Size(assetId string) (width, height int) {
	return internal.AssetSize(assetId)
}
func IsLoaded(assetId string) bool {
	return internal.IsLoaded(assetId)
}

//=================================================================
// private

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

func tryInitShader() {
	if internal.ShaderText.ID == 0 {
		internal.ShaderText = rl.LoadShaderFromMemory("", frag)
	}
}

func getIdPath(p string) string {
	p = text.Replace(p, "\\", "/")
	p = path.RemoveExtension(p)
	return p
}
