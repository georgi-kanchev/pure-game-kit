package assets

import (
	"pure-game-kit/data/path"
	"pure-game-kit/internal"
	"pure-game-kit/utility/text"
	"pure-game-kit/window"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func Size(assetId string) (width, height float32) {
	var w, h = internal.AssetSize(assetId)
	return float32(w), float32(h)
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
	var root = path.Folder(path.Executable()) + path.Divider()
	var result = path.RemoveExtension(text.Remove(p, root))
	result = text.Replace(result, "\\", "/")
	return result
}
