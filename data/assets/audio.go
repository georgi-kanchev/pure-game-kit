package assets

import (
	"maps"
	"pure-game-kit/data/file"
	"pure-game-kit/debug"
	"pure-game-kit/internal"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func LoadSound(filePath string) string {
	tryCreateWindow()
	tryInitAudio()

	var _, has = internal.Sounds[filePath]
	if has {
		return filePath
	}

	if !file.IsExisting(filePath) {
		debug.LogError("Failed to find audio file: \"", filePath, "\"")
		return ""
	}

	var sound = rl.LoadSound(filePath)
	if sound.FrameCount == 0 {
		debug.LogError("Failed to load audio file: \"", filePath, "\"")
		return ""
	}

	internal.Sounds[filePath] = &sound
	return filePath
}
func LoadMusic(filePath string) string {
	tryCreateWindow()
	tryInitAudio()

	var _, has = internal.Music[filePath]
	if has {
		return filePath
	}

	if !file.IsExisting(filePath) {
		debug.LogError("Failed to find audio file: \"", filePath, "\"")
		return ""
	}

	var music = rl.LoadMusicStream(filePath)
	if music.FrameCount == 0 {
		debug.LogError("Failed to load audio file: \"", filePath, "\"")
		return ""
	}

	music.Looping = false
	internal.Music[filePath] = &music
	return filePath
}

func UnloadSound(soundId string) {
	var sound, has = internal.Sounds[soundId]

	if has && !isDefault(soundId) {
		delete(internal.Sounds, soundId)
		rl.UnloadSound(*sound)
	}
}
func UnloadMusic(musicId string) {
	var music, has = internal.Music[musicId]

	if has {
		delete(internal.Music, musicId)
		rl.UnloadMusicStream(*music)
	}
}

func ReloadAllSounds() {
	var loaded = maps.Keys(internal.Sounds)
	UnloadAllSounds()
	for id := range loaded {
		LoadSound(id)
	}
}
func ReloadAllMusic() {
	var loaded = maps.Keys(internal.Music)
	UnloadAllMusic()
	for id := range loaded {
		LoadMusic(id)
	}
}

func UnloadAllSounds() {
	for id := range internal.Sounds {
		UnloadSound(id)
	}
}
func UnloadAllMusic() {
	for id := range internal.Music {
		UnloadMusic(id)
	}
}
