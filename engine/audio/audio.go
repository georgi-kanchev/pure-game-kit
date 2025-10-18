package audio

import (
	"pure-kit/engine/execution/condition"
	"pure-kit/engine/internal"
	"pure-kit/engine/utility/number"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type Audio struct {
	AssetId                  string
	Volume, Pitch, LeftRight float32
	IsLooping, IsPaused      bool

	prevAssetId string
	prevPitch, prevLeftRight,
	prevTime float32
	prevPause bool

	time, duration float32

	// this is field does not contain the sound data
	// it's in internal.Sounds which this instance just uses (raylib alias)
	instance *rl.Sound
}

var Volume, VolumeMusic, VolumeSound float32 = 1, 1, 1

func New(assetId string) *Audio {
	var audio = &Audio{AssetId: assetId, Volume: 1, Pitch: 1, LeftRight: 0.5}
	condition.CallFor(number.ValueMaximum[float32](), audio.update)
	return audio
}

func (audio *Audio) Play() {
	var sound, _ = internal.Sounds[audio.AssetId]
	var music, _ = internal.Music[audio.AssetId]

	audio.IsPaused = false
	if audio.instance != nil && sound != nil {
		rl.PlaySound(*audio.instance)
	} else if music != nil {
		rl.PlayMusicStream(*music)
	}
}
func (audio *Audio) IsPlaying() bool {
	var sound, _ = internal.Sounds[audio.AssetId]
	var music, _ = internal.Music[audio.AssetId]

	if audio.instance != nil && sound != nil {
		return rl.IsSoundPlaying(*audio.instance)
	} else if music != nil {
		return rl.IsMusicStreamPlaying(*music)
	}
	return false
}

func (audio *Audio) Time() (current, duration float32) {
	var sound, _ = internal.Sounds[audio.AssetId]
	var music, _ = internal.Music[audio.AssetId]

	if music != nil {
		return rl.GetMusicTimePlayed(*music), audio.duration // cached upon AssetId change
	}
	if audio.instance != nil && sound != nil {
		return audio.time, audio.duration
	}
	return number.NaN(), number.NaN()
}

//=================================================================
// private

// detects all changes in the values so that the audio can reflect them instantly
func (audio *Audio) update(float32) {
	var sound, _ = internal.Sounds[audio.AssetId]
	var music, _ = internal.Music[audio.AssetId]
	var prevSound, _ = internal.Sounds[audio.prevAssetId]
	var prevMusic, _ = internal.Music[audio.prevAssetId]

	if sound == nil && music == nil {
		if audio.instance != nil { // sound data was unloaded but we're still pointing to it, clear up
			rl.UnloadSoundAlias(*audio.instance)
			audio.instance = nil
		}
		return
	}

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

	// current time cache
	if audio.instance != nil && sound != nil {
		audio.time = currentTime(audio.instance.Stream)
	} else if music != nil {
		audio.time = currentTime(music.Stream)
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

	// volume
	if sound != nil && audio.instance != nil {
		audio.instance.Stream.Buffer.Volume = audio.Volume * VolumeSound * Volume
	}
	if music != nil {
		music.Stream.Buffer.Volume = audio.Volume * VolumeMusic * Volume
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
		rl.SetSoundPan(*audio.instance, audio.LeftRight)
	}
	if music != nil && audio.LeftRight != audio.prevLeftRight {
		rl.SetMusicPan(*music, audio.LeftRight)
	}

	// loop
	if audio.IsLooping && audio.time < audio.prevTime {
		audio.Play()
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

	audio.prevAssetId = audio.AssetId
	audio.prevLeftRight = audio.LeftRight
	audio.prevPitch = audio.Pitch
	audio.prevPause = audio.IsPaused
	audio.prevTime = audio.time
}

func currentTime(stream rl.AudioStream) float32 {
	var processed = stream.Buffer.FramesProcessed
	var cur = stream.Buffer.FrameCursorPos
	var rate = stream.SampleRate
	return float32(processed)/float32(rate) + float32(cur)/float32(rate)
}
