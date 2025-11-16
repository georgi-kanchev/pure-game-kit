package assets

import (
	"pure-game-kit/internal"
	"pure-game-kit/utility/text"
	"pure-game-kit/window"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func UnloadAll() {
	UnloadAllTextures()
	UnloadAllSounds()
	UnloadAllMusic()
	UnloadAllTiledMaps()
	UnloadAllTiledTilesets()
	UnloadAllTiledProjects()
}
func ReloadAll() {
	ReloadAllTextures()
	ReloadAllSounds()
	ReloadAllMusic()
	ReloadAllTiledMaps()
	ReloadAllTiledTilesets()
	ReloadAllTiledProjects()
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
// private

const defaultFont = ""
const defaultTexture = ""
const defaultCursors = "^"
const defaultIcons = "@"
const defaultInputLeft = "["
const defaultInputRight = "]"
const defaultPatterns = "&"
const defaultRetroAtlas = "#"
const defaultUI = "!"
const defaultSoundsUI = "~"

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

func isDefault(id string) bool {
	return !text.Contains(id, ".") // files have '.' in them (folder/name.extension) but default asset ids don't
}
