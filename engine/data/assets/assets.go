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
	var root = path.Folder(path.Executable())
	absolutePath = path.New(root, p)
	id = path.RemoveExtension(p)
	return
}
