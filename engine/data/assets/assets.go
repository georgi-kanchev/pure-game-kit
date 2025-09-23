package assets

import (
	"pure-kit/engine/data/path"
	"pure-kit/engine/internal"
	"pure-kit/engine/window"

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

func getIdPath(p string) (id, absolutePath string) {
	absolutePath = path.New(path.Folder(path.Executable()), p)
	id = path.RemoveExtension(path.LastElement(absolutePath))
	return
}
