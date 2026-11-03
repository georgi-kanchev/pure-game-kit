package audio

import (
	"pure-game-kit/packages/assets"
	"pure-game-kit/packages/internal"
	"pure-game-kit/packages/utility/random"
	"slices"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type Audio struct {
	AssetId                         assets.AudioId
	Volume, Pitch, LeftRightBalance float32

	isPaused  bool
	currentId assets.AudioId // AssetId may be a Variation, this is the raw assetId that is currently playing (or being paused)

	lastPlayFrame uint64 // 2-frame buffer for IsJustFinished()
}

var Volume, VolumeMusic, VolumeSound float32 = 1, 1, 1

func New(assetId assets.AudioId) Audio {
	return Audio{AssetId: assetId, Volume: 1, Pitch: 1, LeftRightBalance: 0.5}
}

//=================================================================

func (a *Audio) Play() {
	a.lastPlayFrame = internal.Frame

	if a.currentId != 0 && a.isPaused {
		var sounds, hasSound = internal.Sounds[int16(a.currentId)]
		var music, hasMusic = internal.Music[int16(a.currentId)]

		a.isPaused = false
		a.ApplyProperties()

		if hasSound {
			for _, s := range sounds {
				rl.ResumeSound(s)
			}
		} else if hasMusic {
			rl.ResumeMusicStream(music)
		}
		return
	}

	var variations, isVariation = internal.AudioVariations[int16(a.AssetId)]
	var assetId = int16(a.AssetId)
	if isVariation {
		assetId = random.Pick(variations...)
	}
	a.isPaused = false
	a.currentId = assets.AudioId(assetId)

	var _, hasSound = internal.Sounds[assetId]
	var music, hasMusic = internal.Music[assetId]
	a.ApplyProperties()
	if hasSound {
		rl.PlaySound(internal.GetSound(assetId, true))
	} else if hasMusic {
		rl.PlayMusicStream(music)
	}
}
func (a *Audio) Pause() {
	if a.isPaused || !a.IsPlaying() {
		return
	}
	a.isPaused = true
	a.lastPlayFrame = internal.Frame

	var sounds, hasSound = internal.Sounds[int16(a.currentId)]
	var music, hasMusic = internal.Music[int16(a.currentId)]
	if hasSound {
		for _, s := range sounds {
			rl.PauseSound(s)
		}
	} else if hasMusic {
		rl.PauseMusicStream(music)
	}
}
func (a *Audio) Stop() {
	if a.currentId == 0 {
		return
	}
	a.isPaused = false
	a.lastPlayFrame = internal.Frame

	var sounds, hasSound = internal.Sounds[int16(a.currentId)]
	var music, hasMusic = internal.Music[int16(a.currentId)]
	if hasSound {
		for _, s := range sounds {
			rl.StopSound(s)
		}
	} else if hasMusic && rl.IsMusicStreamPlaying(music) { // <- prevents raylib stutter quirk in the beginning if Stop() -> Play()
		rl.StopMusicStream(music)
	}
}

func (a *Audio) ApplyProperties() {
	var variations, isVariation = internal.AudioVariations[int16(a.AssetId)]
	if isVariation {
		for _, v := range variations {
			a.applyProperties(v)
		}
		return
	}
	a.applyProperties(int16(a.AssetId))
}

func (a *Audio) IsPaused() bool {
	return a.isPaused
}
func (a *Audio) IsPlaying() bool {
	if a.currentId == 0 {
		return false
	}
	var sounds, _ = internal.Sounds[int16(a.currentId)]
	var music, hasMusic = internal.Music[int16(a.currentId)]
	if hasMusic {
		return rl.IsMusicStreamPlaying(music)
	}
	return slices.ContainsFunc(sounds, rl.IsSoundPlaying)
}
func (a *Audio) IsJustFinished() bool {
	if a.isPaused || a.currentId == 0 {
		return false
	}
	if a.IsPlaying() {
		a.lastPlayFrame = internal.Frame
		return false
	}
	return internal.Frame-1 == a.lastPlayFrame
}

// Returns the currently playing (or paused) AudioId. Useful when Audio.AssetId is a Variation.
func (a *Audio) CurrentAssetId() assets.AudioId {
	return a.currentId
}

// private ========================================================

func (a *Audio) applyProperties(assetId int16) {
	var sounds, hasSound = internal.Sounds[assetId]
	var music, hasMusic = internal.Music[assetId]
	if hasSound {
		for _, sound := range sounds {
			if sound.Stream.Buffer != nil {
				rl.SetSoundVolume(sound, a.Volume*VolumeSound*Volume)
				rl.SetSoundPitch(sound, a.Pitch)
				rl.SetSoundPan(sound, 1-a.LeftRightBalance)
			}
		}
	} else if hasMusic {
		rl.SetMusicVolume(music, a.Volume*VolumeMusic*Volume)
		rl.SetMusicPitch(music, a.Pitch)
		rl.SetMusicPan(music, 1-a.LeftRightBalance)
	}
}
