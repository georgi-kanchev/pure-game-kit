// Loads, reloads or unloads any file data that other engine packages may need.
// Once loaded, an asset id string is provided and the data can be accessed only through it.
// The asset id closely resembles the file path that the loading happened from (not always identical).
//
// Optionally, some embedded default images, sounds and a font may be loaded as a placeholder
// to get things on-screen quickly.
//
// Requires a window because textures & fonts need their OpenGL context,
// and audio needs a device & context to be created.
// Will force initialize a new window whenever needed if there is none.
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
