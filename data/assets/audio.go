package assets

import (
	"maps"
	"pure-game-kit/data/file"
	"pure-game-kit/debug"
	"pure-game-kit/internal"
	"pure-game-kit/utility/collection"
	"slices"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func LoadedSoundIds() []string {
	return collection.MapKeys(internal.Sounds)
}
func LoadedMusicIds() []string {
	return collection.MapKeys(internal.Music)
}

func LoadSound(filePath string) string {
	tryCreateWindow()
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
	var loaded = slices.Collect(maps.Keys(internal.Sounds))
	UnloadAllSounds()
	for _, id := range loaded {
		LoadSound(id)
	}
}
func ReloadAllMusic() {
	var loaded = slices.Collect(maps.Keys(internal.Music))
	UnloadAllMusic()
	for _, id := range loaded {
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
