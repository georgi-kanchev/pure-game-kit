package assets

import (
	"path/filepath"
	"pure-kit/engine/data/file"
	"pure-kit/engine/data/folder"
	"pure-kit/engine/internal"
	"strings"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func LoadSounds(filePaths ...string) []string {
	tryCreateWindow()
	tryInitAudio()

	var result = []string{}
	for _, path := range filePaths {
		var absolutePath = filepath.Join(folder.PathOfExecutable(), path)
		path = strings.ReplaceAll(path, file.Extension(path), "")
		var _, has = internal.Sounds[path]

		if has || !file.Exists(absolutePath) {
			continue
		}

		var sound = rl.LoadSound(absolutePath)

		if sound.FrameCount != 0 {
			internal.Sounds[path] = &sound
			result = append(result, path)
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
