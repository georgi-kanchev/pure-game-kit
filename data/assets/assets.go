package assets

import (
	"pure-game-kit/internal"
	"pure-game-kit/window"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func UnloadAll() {
	UnloadAllTextures()
	UnloadAllSounds()
	UnloadAllMusic()
	UnloadAllTiledMaps()
	UnloadAllTiledTilesets()
}
func ReloadAll() {
	ReloadAllTextures()
	ReloadAllSounds()
	ReloadAllMusic()
	ReloadAllTiledMaps()
	ReloadAllTiledTilesets()
}

//=================================================================

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
