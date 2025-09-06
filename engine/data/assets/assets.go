package assets

import (
	"path/filepath"
	"pure-kit/engine/data/file"
	"pure-kit/engine/data/folder"
	"pure-kit/engine/internal"
	"pure-kit/engine/window"
	"strings"

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

func getIdPath(path string) (id, absolutePath string) {
	absolutePath = filepath.Join(folder.PathOfExecutable(), path)
	id = strings.ReplaceAll(path, file.Extension(path), "")
	return
}
