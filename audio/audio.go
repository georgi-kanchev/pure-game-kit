// A combined controller for the Sound & Music assets. It exists independently of the assets.
package audio

import (
	"pure-game-kit/execution/condition"
	"pure-game-kit/internal"
	"pure-game-kit/utility/number"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type Audio struct {
	AssetId string
	Volume, Pitch, LeftRight,
	FadeIn, FadeOut float32
	IsLooping, IsPaused bool

	prevAssetId string
	prevPitch, prevLeftRight,
	prevTime float32
	prevPause, justFinished bool

	time, duration float32

	// this field does not contain the sound data
	// it's in internal.Sounds which this instance just uses as a raylib alias
	instance rl.Sound
}

var Volume, VolumeMusic, VolumeSound float32 = 1, 1, 1

func New(assetId string) Audio {
	var audio = Audio{AssetId: assetId, Volume: 1, Pitch: 1, LeftRight: 0.5, prevLeftRight: number.NaN()}
	condition.CallFor(number.ValueMaximum[float32](), audio.update)
	return audio
}

//=================================================================

func (a Audio) Play() {
	a.tryReload()

	var _, hasSound = internal.Sounds[a.AssetId]
	var music, hasMusic = internal.Music[a.AssetId]
	var volume = a.volume()

	a.IsPaused = false
	if a.instance.FrameCount != 0 && hasSound {
		rl.PlaySound(a.instance)
		a.instance.Stream.Buffer.Volume = volume
	} else if a.instance.FrameCount != 0 && !hasSound {
		a.instance = rl.Sound{}
	} else if hasMusic {
		rl.PlayMusicStream(music)
		music.Stream.Buffer.Volume = volume
	}
}

//=================================================================

func (a Audio) IsPlaying() bool {
	a.tryReload()

	var _, hasSound = internal.Sounds[a.AssetId]
	var music, hasMusic = internal.Music[a.AssetId]

	if a.instance.FrameCount != 0 && hasSound {
		return rl.IsSoundPlaying(a.instance)
	} else if a.instance.FrameCount != 0 && !hasSound {
		a.instance = rl.Sound{}
	} else if hasMusic {
		return rl.IsMusicStreamPlaying(music)
	}
	return false
}
func (a Audio) IsJustFinished() bool {
	return a.justFinished
}

func (a Audio) Time() (current, duration float32) {
	return a.time, a.duration
}

//=================================================================
// private

// detects all changes in the values so that the audio can reflect them instantly
func (a Audio) update(float32) {
	var _, hasSound = internal.Sounds[a.AssetId]
	var music, hasMusic = internal.Music[a.AssetId]

	if !hasSound && !hasMusic {
		if a.instance.FrameCount != 0 { // sound data was unloaded but we're still pointing to it, clear up
			rl.UnloadSoundAlias(a.instance)
			a.instance = rl.Sound{}
		}
		return
	}

	a.tryReload()

	// current time cache
	if a.instance.FrameCount != 0 && hasSound {
		a.time = currentTime(a.instance.Stream)
	} else if hasMusic {
		a.time = rl.GetMusicTimePlayed(music)
	}

	// volume
	var volume = a.volume()
	if hasSound && a.instance.FrameCount != 0 {
		a.instance.Stream.Buffer.Volume = volume
	}
	if hasMusic {
		music.Stream.Buffer.Volume = volume
	}

	// pitch
	if hasSound && a.instance.FrameCount != 0 && a.Pitch != a.prevPitch {
		rl.SetSoundPitch(a.instance, a.Pitch)
	}
	if hasMusic && a.Pitch != a.prevPitch {
		rl.SetMusicPitch(music, a.Pitch)
	}

	// leftRight
	if hasSound && a.instance.FrameCount != 0 && a.LeftRight != a.prevLeftRight {
		rl.SetSoundPan(a.instance, 1-a.LeftRight)
	}
	if hasMusic && a.LeftRight != a.prevLeftRight {
		rl.SetMusicPan(music, 1-a.LeftRight)
	}

	// loop
	a.justFinished = false
	if a.time < a.prevTime {
		if a.IsLooping {
			a.Play()
		} else {
			a.justFinished = true
		}
	}

	// pause
	if hasSound && a.instance.FrameCount != 0 && a.IsPaused != a.prevPause {
		if a.IsPaused {
			rl.PauseSound(a.instance)
		} else {
			rl.ResumeSound(a.instance)
		}
	}
	if hasMusic && a.IsPaused != a.prevPause {
		if a.IsPaused {
			rl.PauseMusicStream(music)
		} else {
			rl.ResumeMusicStream(music)
		}
	}

	a.prevLeftRight = a.LeftRight
	a.prevPitch = a.Pitch
	a.prevPause = a.IsPaused
	a.prevTime = a.time
}

func (a Audio) volume() float32 {
	var fadeIn = number.Limit(number.Map(a.time, 0, a.FadeIn, 0, 1), 0, 1)
	var fadeOut = number.Limit(number.Map(a.time, a.duration-a.FadeOut, a.duration, 1, 0), 0, 1)
	if a.FadeIn <= 0 {
		fadeIn = 1
	}
	if a.FadeOut <= 0 {
		fadeOut = 1
	}

	return a.Volume * VolumeSound * Volume * fadeIn * fadeOut
}
func currentTime(stream rl.AudioStream) float32 {
	var processed = stream.Buffer.FramesProcessed
	var cur = stream.Buffer.FrameCursorPos
	var rate = stream.SampleRate
	return float32(processed)/float32(rate) + float32(cur)/float32(rate)
}

func (a Audio) tryReload() {
	var sound, hasSound = internal.Sounds[a.AssetId]
	var music, hasMusic = internal.Music[a.AssetId]
	var prevSound, hasPrevSound = internal.Sounds[a.prevAssetId]
	var prevMusic, hasPrevMusic = internal.Music[a.prevAssetId]

	// stop & cleanup
	if sound != prevSound {
		a.duration = 0
		if hasPrevSound && a.instance.FrameCount != 0 {
			rl.UnloadSoundAlias(a.instance)
		}
	}
	if music != prevMusic {
		a.duration = 0
		if hasPrevMusic {
			rl.StopMusicStream(prevMusic)
		}
	}
	// assetId
	if hasSound && a.AssetId != a.prevAssetId {
		var newInstance = rl.LoadSoundAlias(sound)
		a.instance = newInstance // load & set the new instance
		a.time = 0
		a.duration = float32(sound.FrameCount) / float32(sound.Stream.SampleRate)
	}
	if hasMusic && a.AssetId != a.prevAssetId {
		a.time = 0
		a.duration = rl.GetMusicTimeLength(music)
	}

	a.prevAssetId = a.AssetId
}
