package internal

import rl "github.com/gen2brain/raylib-go/raylib"

var Sounds = make(map[int16][]rl.Sound) // [0]soundData; [1].FrameCount = playIndex; [2...]soundAliases
var Music = make(map[int16]rl.Music)
var AudioUpdates []func()

func GetSound(id int16, prepareNext bool) rl.Sound {
	var sounds = Sounds[id]
	if len(sounds) == 1 {
		return sounds[0]
	}
	var index = sounds[1].FrameCount + 2
	var result = sounds[index]
	if prepareNext {
		sounds[1].FrameCount = uint32(int(sounds[1].FrameCount+1) % (len(sounds) - 2))
	}
	return result
}

func UpdateAudio() {
	for _, m := range Music {
		rl.UpdateMusicStream(m)
	}
	for _, u := range AudioUpdates {
		u()
	}
}
