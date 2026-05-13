// A combined controller for the Sound & Music assets. It exists independently of the assets.
package audio

import (
	"pure-game-kit/packages/assets"
	"pure-game-kit/packages/internal"
	"pure-game-kit/packages/utility/number"
	"slices"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type Audio struct {
	AssetId assets.AudioId
	Volume, Pitch, LeftRightBalance,
	FadeIn, FadeOut float32
	IsLooping, IsPaused bool

	//=================================================================

	prevAssetId              string
	prevPitch, prevLeftRight float32
	prevPause, justFinished  bool

	lastPlayTime, duration float32
	finishFlag             bool
}

var Volume, VolumeMusic, VolumeSound float32 = 1, 1, 1

func New(assetId assets.AudioId) *Audio {
	var result = &Audio{AssetId: assetId, Volume: 1, Pitch: 1, LeftRightBalance: 0.5}
	internal.AudioUpdates = append(internal.AudioUpdates, func() {
		result.update()
	})
	return result
}

//=================================================================

func (a *Audio) Play() {
	var _, hasSound = internal.Sounds[int16(a.AssetId)]
	var music, hasMusic = internal.Music[int16(a.AssetId)]
	var volume = a.volume()

	a.lastPlayTime = internal.Runtime
	a.IsPaused = false
	if hasSound {
		var sound = internal.GetSound(int16(a.AssetId), true)
		rl.PlaySound(sound)
		sound.Stream.Buffer.Volume = volume
	} else if hasMusic {
		rl.PlayMusicStream(music)
		music.Stream.Buffer.Volume = volume
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
	return a.justFinished
}

// private ========================================================

// detects all changes in the values so that the audio can reflect them instantly
func (a *Audio) update() {
	var _, hasSound = internal.Sounds[int16(a.AssetId)]
	var music, hasMusic = internal.Music[int16(a.AssetId)]
	var volume = a.volume()

	if hasSound {
		var sound = internal.GetSound(int16(a.AssetId), false)
		sound.Stream.Buffer.Volume = volume

		if a.Pitch != a.prevPitch {
			rl.SetSoundPitch(sound, a.Pitch)
		}
		if a.LeftRightBalance != a.prevLeftRight {
			rl.SetSoundPan(sound, 1-a.LeftRightBalance)
		}
		if a.IsPaused != a.prevPause {
			if a.IsPaused {
				rl.PauseSound(sound)
			} else {
				rl.ResumeSound(sound)
			}
		}
		if a.duration == 0 {
			a.duration = float32(sound.FrameCount) / float32(sound.Stream.SampleRate)
		}
	} else if hasMusic {
		music.Stream.Buffer.Volume = volume

		if a.Pitch != a.prevPitch {
			rl.SetMusicPitch(music, a.Pitch)
		}
		if a.LeftRightBalance != a.prevLeftRight {
			rl.SetMusicPan(music, 1-a.LeftRightBalance)
		}
		if a.IsPaused != a.prevPause {
			if a.IsPaused {
				rl.PauseMusicStream(music)
			} else {
				rl.ResumeMusicStream(music)
			}
		}
		if a.duration == 0 {
			a.duration = rl.GetMusicTimeLength(music)
		}
	}

	a.justFinished = false
	if internal.Runtime-a.lastPlayTime > a.duration {
		if a.IsLooping {
			a.Play()
		} else {
			a.justFinished = true
		}
	}

	a.prevLeftRight = a.LeftRightBalance
	a.prevPitch = a.Pitch
	a.prevPause = a.IsPaused
}

func (a *Audio) volume() float32 {
	var progress = number.Limit(number.Map(internal.Runtime, a.lastPlayTime, a.lastPlayTime+a.duration, 0, 1), 0, 1)
	var fadeIn = number.Limit(number.Map(progress, 0, a.FadeIn, 0, 1), 0, 1)
	var fadeOut = number.Limit(number.Map(progress, a.duration-a.FadeOut, a.duration, 1, 0), 0, 1)
	if a.FadeIn <= 0 {
		fadeIn = 1
	}
	if a.FadeOut <= 0 {
		fadeOut = 1
	}

	return a.Volume * VolumeSound * Volume * fadeIn * fadeOut
}
