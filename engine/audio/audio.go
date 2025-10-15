package audio

import (
	"pure-kit/engine/internal"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type Audio struct {
	AssetId                  string
	Volume, Pitch, LeftRight float32
}

var GlobalVolume float32 = 1

func New(assetId string) *Audio {
	return &Audio{AssetId: assetId, Volume: 1, Pitch: 1, LeftRight: 0.5}
}

func (audio *Audio) Play() {
	var sounds, isSound = internal.Sounds[audio.AssetId]
	var music, _ = internal.Music[audio.AssetId]

	if isSound {
		for _, s := range sounds {
			if !rl.IsSoundPlaying(*s) {
				rl.SetSoundPitch(*s, audio.Pitch)
				rl.SetSoundVolume(*s, audio.Volume*GlobalVolume)
				rl.SetSoundPan(*s, 1-audio.LeftRight)
				rl.PlaySound(*s)
				break
			}
		}
	} else if music != nil {
		rl.SetMusicPitch(*music, audio.Pitch)
		rl.SetMusicVolume(*music, audio.Volume*GlobalVolume)
		rl.SetMusicPan(*music, 1-audio.LeftRight)
		rl.PlayMusicStream(*music)
	}
}

func (audio *Audio) IsPlaying() bool {
	var sounds, isSound = internal.Sounds[audio.AssetId]
	var music, _ = internal.Music[audio.AssetId]

	if isSound {
		for _, s := range sounds {
			if rl.IsSoundPlaying(*s) {
				return true
			}
		}
	} else if music != nil {
		return rl.IsMusicStreamPlaying(*music)
	}
	return false
}
