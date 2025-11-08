package assets

import (
	"pure-game-kit/data/file"
	"pure-game-kit/debug"
	"pure-game-kit/internal"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func LoadSound(filePath string) string {
	tryCreateWindow()
	tryInitAudio()

	var result = ""
	var id = getIdPath(filePath)

	if !file.IsExisting(filePath) {
		debug.LogError("Failed to find audio file: \"", filePath, "\"")
		return result
	}

	var sound = rl.LoadSound(filePath)
	if sound.FrameCount != 0 {
		internal.Sounds[id] = &sound
		result = id
	} else {
		debug.LogError("Failed to load audio file: \"", filePath, "\"")
	}

	return result
}
func LoadMusic(filePath string) string {
	tryCreateWindow()
	tryInitAudio()

	var result = ""
	var id = getIdPath(filePath)

	if !file.IsExisting(filePath) {
		debug.LogError("Failed to find audio file: \"", filePath, "\"")
		return result
	}

	var music = rl.LoadMusicStream(filePath)
	if music.FrameCount != 0 {
		music.Looping = false
		internal.Music[id] = &music
		result = id
	} else {
		debug.LogError("Failed to load audio file: \"", filePath, "\"")
	}

	return result
}

func ReloadAllSounds() {
	for id := range internal.Sounds {
		LoadSound(id)
	}
}
func ReloadAllMusic() {
	for id := range internal.Music {
		LoadMusic(id)
	}
}

func UnloadSound(soundId string) {
	var sound, has = internal.Sounds[soundId]

	if has {
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
