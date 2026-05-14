package audio

import (
	"math" // Added for math.Round
	"pure-game-kit/packages/assets"
	"pure-game-kit/packages/internal"
	"pure-game-kit/packages/utility/number"
	"slices"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type Audio struct {
	AssetId                         assets.AudioId
	Volume, Pitch, LeftRightBalance float32

	playTick, pauseTick uint64
	isPaused            bool
}

var Volume, VolumeMusic, VolumeSound float32 = 1, 1, 1

func New(assetId assets.AudioId) Audio {
	return Audio{AssetId: assetId, Volume: 1, Pitch: 1, LeftRightBalance: 0.5, playTick: 9999999999999, isPaused: false}
}

//=================================================================

func (a *Audio) Play() {
	a.isPaused = false

	var _, hasSound = internal.Sounds[int16(a.AssetId)]
	var music, hasMusic = internal.Music[int16(a.AssetId)]

	a.playTick = internal.Tick
	if hasSound {
		var sound = internal.GetSound(int16(a.AssetId), true)
		a.ApplyProperties()
		rl.PlaySound(sound)
	} else if hasMusic {
		a.ApplyProperties()
		rl.PlayMusicStream(music)
	}
}

func (a *Audio) Pause() {
	if a.isPaused {
		return
	}

	a.isPaused = true
	a.pauseTick = internal.Tick

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
func (a *Audio) Resume() {
	if !a.isPaused {
		return
	}

	a.isPaused = false
	a.playTick += internal.Tick - a.pauseTick // shift by ticks spent paused

	var sounds, hasSound = internal.Sounds[int16(a.AssetId)]
	var music, hasMusic = internal.Music[int16(a.AssetId)]

	a.ApplyProperties()
	if hasSound {
		for _, s := range sounds {
			rl.ResumeSound(s)
		}
	} else if hasMusic {
		rl.ResumeMusicStream(music)
	}
}

func (a *Audio) ApplyProperties() {
	var sounds, hasSound = internal.Sounds[int16(a.AssetId)]
	var music, hasMusic = internal.Music[int16(a.AssetId)]
	var volume = a.Volume * Volume

	if hasSound {
		for _, sound := range sounds {
			if sound.Stream.Buffer != nil {
				sound.Stream.Buffer.Volume = volume * VolumeSound
				rl.SetSoundPitch(sound, a.Pitch)
				rl.SetSoundPan(sound, 1-a.LeftRightBalance)
			}
		}
	} else if hasMusic {
		music.Stream.Buffer.Volume = volume * VolumeMusic
		rl.SetMusicPitch(music, a.Pitch)
		rl.SetMusicPan(music, 1-a.LeftRightBalance)
	}
}

//=================================================================

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
	var durationInTicks uint64
	var tps = float32(internal.TargetTPS)

	if hasSound {
		for _, s := range sounds {
			durationInTicks = uint64(math.Round(float64(float32(s.FrameCount) / float32(s.Stream.SampleRate) * tps)))
		}
		var finishTick = a.playTick + uint64(math.Round(float64(float32(durationInTicks)/a.Pitch)))
		return internal.Tick == finishTick

	} else if hasMusic {
		durationInTicks = uint64(math.Round(float64(rl.GetMusicTimeLength(music) * tps)))
		var finishTick = a.playTick + uint64(math.Round(float64(float32(durationInTicks)/a.Pitch)))
		return internal.Tick == finishTick
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
