package assets

import (
	"pure-kit/engine/data/file"
	"pure-kit/engine/internal"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func LoadSounds(filePaths ...string) []string {
	tryCreateWindow()
	tryInitAudio()

	var result = []string{}
	for _, path := range filePaths {
		var id, absolutePath = getIdPath(path)
		var _, has = internal.Sounds[id]

		if has || !file.Exists(absolutePath) {
			continue
		}

		var sound = rl.LoadSound(absolutePath)

		if sound.FrameCount != 0 {
			internal.Sounds[id] = &sound
			result = append(result, id)
		}
	}

	return result
}

func UnloadSounds(soundIds ...string) {
	for _, v := range soundIds {
		var sound, has = internal.Sounds[v]

		if has {
			delete(internal.Sounds, v)
			rl.UnloadSound(*sound)
		}
	}
}
