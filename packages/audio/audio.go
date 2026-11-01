package audio

import (
	"pure-game-kit/packages/assets"
	"pure-game-kit/packages/internal"
	"pure-game-kit/packages/utility/number"
	"slices"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type Audio struct {
	AssetId                         assets.AudioId
	Volume, Pitch, LeftRightBalance float32

	playTime, pauseTime, finishTime float32
	finishedReported, isPaused      bool
}

var Volume, VolumeMusic, VolumeSound float32 = 1, 1, 1

func New(assetId assets.AudioId) Audio {
	return Audio{AssetId: assetId, Volume: 1, Pitch: 1, LeftRightBalance: 0.5, playTime: number.NaN()}
}

//=================================================================

func (a *Audio) Play() {
	var sounds, hasSound = internal.Sounds[int16(a.AssetId)]
	var music, hasMusic = internal.Music[int16(a.AssetId)]

	a.finishedReported = false
	a.ApplyProperties()

	if a.isPaused {
		a.isPaused = false
		a.playTime += internal.Runtime - a.pauseTime // shift by time spent paused

		if hasSound {
			for _, s := range sounds {
				rl.ResumeSound(s)
			}
		} else if hasMusic {
			rl.ResumeMusicStream(music)
		}
		return
	}

	a.playTime = internal.Runtime
	if hasSound {
		rl.PlaySound(internal.GetSound(int16(a.AssetId), true))
	} else if hasMusic {
		rl.PlayMusicStream(music)
	}
}
func (a *Audio) Pause() {
	if a.isPaused {
		return
	}

	a.isPaused = true
	a.pauseTime = internal.Runtime

	var sounds, hasSound = internal.Sounds[int16(a.AssetId)]
	var music, hasMusic = internal.Music[int16(a.AssetId)]
	if hasSound {
		for _, s := range sounds {
			rl.PauseSound(s)
		}
	} else if hasMusic {
		rl.PauseMusicStream(music)
	}
}
func (a *Audio) Stop() {
	a.isPaused = true
	a.finishedReported = false

	var _, hasSound = internal.Sounds[int16(a.AssetId)]
	var music, hasMusic = internal.Music[int16(a.AssetId)]
	if hasSound {
		rl.StopSound(internal.GetSound(int16(a.AssetId), false))
	} else if hasMusic {
		rl.StopMusicStream(music)
	}
}

func (a *Audio) ApplyProperties() {
	var sounds, hasSound = internal.Sounds[int16(a.AssetId)]
	var music, hasMusic = internal.Music[int16(a.AssetId)]
	if hasSound {
		for _, sound := range sounds {
			if sound.Stream.Buffer != nil {
				sound.Stream.Buffer.Volume = a.Volume * VolumeSound * Volume
				rl.SetSoundPitch(sound, a.Pitch)
				rl.SetSoundPan(sound, 1-a.LeftRightBalance)
			}
		}
	} else if hasMusic {
		music.Stream.Buffer.Volume = a.Volume * VolumeMusic * Volume
		rl.SetMusicPitch(music, a.Pitch)
		rl.SetMusicPan(music, 1-a.LeftRightBalance)
	}
}

func (a *Audio) IsPlaying() bool {
	var sounds, _ = internal.Sounds[int16(a.AssetId)]
	var music, hasMusic = internal.Music[int16(a.AssetId)]

	if hasMusic {
		return rl.IsMusicStreamPlaying(music)
	}
	return slices.ContainsFunc(sounds, rl.IsSoundPlaying)
}
func (a *Audio) IsJustFinished() bool {
	if a.isPaused {
		return false
	}

	var sounds, hasSound = internal.Sounds[int16(a.AssetId)]
	var music, hasMusic = internal.Music[int16(a.AssetId)]
	if !a.finishedReported {
		var duration float32
		if hasSound {
			for _, s := range sounds {
				duration = float32(s.FrameCount) / float32(s.Stream.SampleRate)
			}
		} else if hasMusic {
			duration = rl.GetMusicTimeLength(music)
		}

		a.finishTime = a.playTime + duration/a.Pitch
		if internal.Runtime >= a.finishTime {
			a.finishedReported = true
			return true
		}
	}
	return false
}

// private ========================================================

func currentTime(stream rl.AudioStream) float32 {
	if stream.Buffer == nil {
		return number.NaN()
	}
	var processed = stream.Buffer.FramesProcessed
	var cur = stream.Buffer.FrameCursorPos
	var rate = stream.SampleRate
	return float32(processed)/float32(rate) + float32(cur)/float32(rate)
}
