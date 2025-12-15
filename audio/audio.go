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
	instance *rl.Sound
}

var Volume, VolumeMusic, VolumeSound float32 = 1, 1, 1

func New(assetId string) *Audio {
	var audio = &Audio{AssetId: assetId, Volume: 1, Pitch: 1, LeftRight: 0.5, prevLeftRight: number.NaN()}
	condition.CallFor(number.ValueMaximum[float32](), audio.update)
	return audio
}

//=================================================================

func (audio *Audio) Play() {
	audio.tryReload()

	var sound, _ = internal.Sounds[audio.AssetId]
	var music, _ = internal.Music[audio.AssetId]
	var volume = audio.volume()

	audio.IsPaused = false
	if audio.instance != nil && sound != nil {
		rl.PlaySound(*audio.instance)
		audio.instance.Stream.Buffer.Volume = volume
	} else if audio.instance != nil && sound == nil {
		audio.instance = nil
	} else if music != nil {
		rl.PlayMusicStream(*music)
		music.Stream.Buffer.Volume = volume
	}
}

//=================================================================

func (audio *Audio) IsPlaying() bool {
	audio.tryReload()

	var sound, _ = internal.Sounds[audio.AssetId]
	var music, _ = internal.Music[audio.AssetId]

	if audio.instance != nil && sound != nil {
		return rl.IsSoundPlaying(*audio.instance)
	} else if audio.instance != nil && sound == nil {
		audio.instance = nil
	} else if music != nil {
		return rl.IsMusicStreamPlaying(*music)
	}
	return false
}
func (audio *Audio) IsJustFinished() bool {
	return audio.justFinished
}

func (audio *Audio) Time() (current, duration float32) {
	return audio.time, audio.duration
}

//=================================================================
// private

// detects all changes in the values so that the audio can reflect them instantly
func (audio *Audio) update(float32) {
	var sound, _ = internal.Sounds[audio.AssetId]
	var music, _ = internal.Music[audio.AssetId]

	if sound == nil && music == nil {
		if audio.instance != nil { // sound data was unloaded but we're still pointing to it, clear up
			rl.UnloadSoundAlias(*audio.instance)
			audio.instance = nil
		}
		return
	}

	audio.tryReload()

	// current time cache
	if audio.instance != nil && sound != nil {
		audio.time = currentTime(audio.instance.Stream)
	} else if music != nil {
		audio.time = rl.GetMusicTimePlayed(*music)
	}

	// volume
	var volume = audio.volume()
	if sound != nil && audio.instance != nil {
		audio.instance.Stream.Buffer.Volume = volume
	}
	if music != nil {
		music.Stream.Buffer.Volume = volume
	}

	// pitch
	if sound != nil && audio.instance != nil && audio.Pitch != audio.prevPitch {
		rl.SetSoundPitch(*audio.instance, audio.Pitch)
	}
	if music != nil && audio.Pitch != audio.prevPitch {
		rl.SetMusicPitch(*music, audio.Pitch)
	}

	// leftRight
	if sound != nil && audio.instance != nil && audio.LeftRight != audio.prevLeftRight {
		rl.SetSoundPan(*audio.instance, 1-audio.LeftRight)
	}
	if music != nil && audio.LeftRight != audio.prevLeftRight {
		rl.SetMusicPan(*music, 1-audio.LeftRight)
	}

	// loop
	audio.justFinished = false
	if audio.time < audio.prevTime {
		if audio.IsLooping {
			audio.Play()
		} else {
			audio.justFinished = true
		}
	}

	// pause
	if sound != nil && audio.instance != nil && audio.IsPaused != audio.prevPause {
		if audio.IsPaused {
			rl.PauseSound(*audio.instance)
		} else {
			rl.ResumeSound(*audio.instance)
		}
	}
	if music != nil && audio.IsPaused != audio.prevPause {
		if audio.IsPaused {
			rl.PauseMusicStream(*music)
		} else {
			rl.ResumeMusicStream(*music)
		}
	}

	audio.prevLeftRight = audio.LeftRight
	audio.prevPitch = audio.Pitch
	audio.prevPause = audio.IsPaused
	audio.prevTime = audio.time
}

func (audio *Audio) volume() float32 {
	var fadeIn = number.Limit(number.Map(audio.time, 0, audio.FadeIn, 0, 1), 0, 1)
	var fadeOut = number.Limit(number.Map(audio.time, audio.duration-audio.FadeOut, audio.duration, 1, 0), 0, 1)
	if audio.FadeIn <= 0 {
		fadeIn = 1
	}
	if audio.FadeOut <= 0 {
		fadeOut = 1
	}

	return audio.Volume * VolumeSound * Volume * fadeIn * fadeOut
}

func currentTime(stream rl.AudioStream) float32 {
	var processed = stream.Buffer.FramesProcessed
	var cur = stream.Buffer.FrameCursorPos
	var rate = stream.SampleRate
	return float32(processed)/float32(rate) + float32(cur)/float32(rate)
}

func (audio *Audio) tryReload() {
	var sound, _ = internal.Sounds[audio.AssetId]
	var music, _ = internal.Music[audio.AssetId]
	var prevSound, _ = internal.Sounds[audio.prevAssetId]
	var prevMusic, _ = internal.Music[audio.prevAssetId]

	// stop & cleanup
	if sound != prevSound {
		audio.duration = 0
		if prevSound != nil && audio.instance != nil {
			rl.UnloadSoundAlias(*audio.instance)
		}
	}
	if music != prevMusic {
		audio.duration = 0
		if prevMusic != nil {
			rl.StopMusicStream(*prevMusic)
		}
	}
	// assetId
	if sound != nil && audio.AssetId != audio.prevAssetId {
		var newInstance = rl.LoadSoundAlias(*sound)
		audio.instance = &newInstance // load & set the new instance
		audio.time = 0
		audio.duration = float32(sound.FrameCount) / float32(sound.Stream.SampleRate)
	}
	if music != nil && audio.AssetId != audio.prevAssetId {
		audio.time = 0
		audio.duration = rl.GetMusicTimeLength(*music)
	}

	audio.prevAssetId = audio.AssetId
}
