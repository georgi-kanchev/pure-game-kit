package assets

import (
	"pure-kit/engine/data/file"
	"pure-kit/engine/internal"
	"pure-kit/engine/utility/number"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func LoadSound(filePath string, maximumStacks int) string {
	filePath = internal.MakeAbsolutePath(filePath)
	tryCreateWindow()
	tryInitAudio()

	var result = ""
	var id = getIdPath(filePath)
	var _, has = internal.Sounds[id]

	if has || !file.IsExisting(filePath) {
		return result
	}

	maximumStacks = number.Biggest(1, maximumStacks)
	var sound = rl.LoadSound(filePath)
	if sound.FrameCount != 0 {
		var instances = make([]*rl.Sound, maximumStacks)
		instances[0] = &sound

		for i := 1; i < len(instances); i++ {
			var alias = rl.LoadSoundAlias(sound)
			instances[i] = &alias
		}

		internal.Sounds[id] = instances
		result = id
	}

	return result
}
func LoadMusic(filePath string) string {
	filePath = internal.MakeAbsolutePath(filePath)
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
	var sounds, has = internal.Sounds[soundId]

	if has {
		delete(internal.Sounds, soundId)
		rl.UnloadSound(*sounds[0])
		for i := 1; i < len(sounds); i++ {
			rl.UnloadSoundAlias(*sounds[i])
		}
	}
}
func UnloadMusic(musicId string) {
	var music, has = internal.Music[musicId]

	if has {
		delete(internal.Music, musicId)
		rl.UnloadMusicStream(*music)
	}
}
