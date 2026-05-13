package assets

import (
	"pure-game-kit/packages/internal"
	"pure-game-kit/packages/utility/debug"
	"pure-game-kit/packages/utility/file"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type AudioId int16

func LoadSound(filePath string, maxOverlapCount byte) AudioId {
	tryInitAudio()

	if !file.Exists(filePath) {
		debug.LogError("Failed to find audio file: \"", filePath, "\"")
		return 0
	}

	var sound = rl.LoadSound(filePath)
	if sound.FrameCount == 0 {
		debug.LogError("Failed to load audio file: \"", filePath, "\"")
		return 0
	}

	var id = len(internal.Sounds)
	var result = make([]rl.Sound, 0, maxOverlapCount)
	result = append(result, sound)
	if maxOverlapCount != 0 {
		result = append(result, rl.Sound{}) // .FrameCount used for keeping nextAliasPlayIndex
	}
	for i := 2; i < int(maxOverlapCount+2); i++ {
		result = append(result, rl.LoadSoundAlias(sound))
	}
	internal.Sounds[int16(id)] = result
	return AudioId(id)
}
func LoadMusic(filePath string) AudioId {
	tryInitAudio()

	if !file.Exists(filePath) {
		debug.LogError("Failed to find audio file: \"", filePath, "\"")
		return 0
	}

	var music = rl.LoadMusicStream(filePath)
	if music.FrameCount == 0 {
		debug.LogError("Failed to load audio file: \"", filePath, "\"")
		return 0
	}

	music.Looping = false
	var id = len(internal.Music)
	internal.Music[int16(id)] = music
	return AudioId(id)
}

func (a AudioId) Duration() float32 {
	if a < 0 { // is music
		var music = internal.Music[int16(a)]
		return float32(music.FrameCount) / float32(music.Stream.SampleRate)
	}

	var sound = internal.Sounds[int16(a)][0]
	return float32(sound.FrameCount) / float32(sound.Stream.SampleRate)
}
