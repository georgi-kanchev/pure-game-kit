package assets

import (
	"pure-game-kit/packages/internal"
	"pure-game-kit/packages/utility/debug"
	"pure-game-kit/packages/utility/file"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func LoadSound(filePath string) string {
	tryInitAudio()

	var _, has = internal.Sounds[filePath]
	if has {
		return filePath
	}

	if !file.Exists(filePath) {
		debug.LogError("Failed to find audio file: \"", filePath, "\"")
		return ""
	}

	var sound = rl.LoadSound(filePath)
	if sound.FrameCount == 0 {
		debug.LogError("Failed to load audio file: \"", filePath, "\"")
		return ""
	}

	internal.Sounds[filePath] = sound
	return filePath
}
func LoadMusic(filePath string) string {
	tryInitAudio()

	var _, has = internal.Music[filePath]
	if has {
		return filePath
	}

	if !file.Exists(filePath) {
		debug.LogError("Failed to find audio file: \"", filePath, "\"")
		return ""
	}

	var music = rl.LoadMusicStream(filePath)
	if music.FrameCount == 0 {
		debug.LogError("Failed to load audio file: \"", filePath, "\"")
		return ""
	}

	music.Looping = false
	internal.Music[filePath] = music
	return filePath
}
