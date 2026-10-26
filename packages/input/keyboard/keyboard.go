package keyboard

import (
	"pure-game-kit/packages/input/keyboard/key"
	k "pure-game-kit/packages/input/keyboard/key"
	"pure-game-kit/packages/internal"
	i "pure-game-kit/packages/internal"
)

func Input() []rune { return i.Input }

func IsKeyHeld(key k.Enum, delay float32) bool {
	if holdingKey != k.NONE && holdingKey != key {
		return false // only one key can be held at a time (first one)
	}
	if IsKeyPressed(key) {
		holdingKey, holdDuration = key, holdDuration+internal.FrameDelta
	}
	if IsKeyJustReleased(key) {
		holdDuration, holdingKey = 0, k.NONE
	}
	if !i.Keys[key] || holdDuration < delay {
		return false
	}
	if internal.Runtime > holdTimeStart+0.05 {
		holdTimeStart = internal.Runtime
		return true
	}
	return false
}
func IsKeyPressed(key k.Enum) bool      { return i.Keys[key] }
func IsKeyJustPressed(key k.Enum) bool  { return i.Keys[key] && !i.KeysPrev[key] }
func IsKeyJustReleased(key k.Enum) bool { return !i.Keys[key] && i.KeysPrev[key] }

func IsAnyKeyPressed() bool      { return i.AnyKey }
func IsAnyKeyJustPressed() bool  { return i.AnyKey && !i.AnyKeyPrev }
func IsAnyKeyJustReleased() bool { return !i.AnyKey && i.AnyKeyPrev }

func IsComboJustPressed(keys ...k.Enum) bool {
	return combo(keys) && IsKeyJustPressed(keys[len(keys)-1])
}
func IsComboHeld(delay float32, keys ...k.Enum) bool {
	return combo(keys) && IsKeyHeld(keys[len(keys)-1], delay)
}

//=================================================================

var holdingKey key.Enum = k.NONE
var holdTimeStart, holdDuration float32

func combo(keys []k.Enum) bool {
	if i.KeyCount != len(keys) {
		return false // pressed key count doesn't match the combo, exit early
	}

	for _, kk := range keys {
		if kk == k.NONE || int(kk) >= len(i.Keys) || !i.Keys[kk] {
			return false
		}
	}
	return true
}
