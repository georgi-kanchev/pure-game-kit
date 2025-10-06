package assets

import (
	"pure-kit/engine/data/file"
	"pure-kit/engine/internal"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func LoadSounds(filePath string) []string {
	filePath = internal.MakeAbsolutePath(filePath)
	tryCreateWindow()
	tryInitAudio()

	var result = []string{}
	var id = getIdPath(filePath)
	var _, has = internal.Sounds[id]

	if has || !file.IsExisting(filePath) {
		return result
	}

	var sound = rl.LoadSound(filePath)

	if sound.FrameCount != 0 {
		internal.Sounds[id] = &sound
		result = append(result, id)
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
