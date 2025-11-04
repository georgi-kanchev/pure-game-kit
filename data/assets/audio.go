package assets

import (
	"pure-game-kit/data/file"
	"pure-game-kit/internal"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func LoadSound(filePath string) string {
	tryCreateWindow()
	tryInitAudio()

	var result = ""
	var id = getIdPath(filePath)
	var _, has = internal.Sounds[id]

	if has || !file.IsExisting(filePath) {
		return result
	}

	var sound = rl.LoadSound(filePath)
	if sound.FrameCount != 0 {
		internal.Sounds[id] = &sound
		result = id
	}

	return result
}
func LoadMusic(filePath string) string {
	tryCreateWindow()
	tryInitAudio()

	var result = ""
	var id = getIdPath(filePath)
	var _, has = internal.Music[id]

	if has || !file.IsExisting(filePath) {
		return result
	}

	var music = rl.LoadMusicStream(filePath)
	if music.FrameCount != 0 {
		music.Looping = false
		internal.Music[id] = &music
		result = id
	}

	return result
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
