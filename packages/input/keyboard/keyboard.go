// Contains checks for whether a certain keyboard key is interacted with in various ways.
// Also provides text input and the currently pressed keys. Meant to be checked every frame
// instead of subscribtion-based events/callbacks.
package keyboard

import (
	"pure-game-kit/packages/internal"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func Input() string {
	return internal.Input
}

func IsKeyPressed(key int) bool {
	if key < 0 || key >= len(internal.Keys) {
		return false
	}
	return internal.Keys[key]
}
func IsKeyHeld(key int) bool {
	return rl.IsKeyPressedRepeat(int32(key))
}

func IsKeyJustPressed(key int) bool {
	if key < 0 || key >= len(internal.Keys) {
		return false
	}
	return internal.Keys[key] && !internal.KeysPrev[key]
}
func IsKeyJustReleased(key int) bool {
	if key < 0 || key >= len(internal.Keys) {
		return false
	}
	return !internal.Keys[key] && internal.KeysPrev[key]
}

func IsAnyKeyPressed() bool {
	return internal.AnyKey
}
func IsAnyKeyJustPressed() bool {
	return internal.AnyKey && !internal.AnyKeyPrev
}
func IsAnyKeyJustReleased() bool {
	return !internal.AnyKey && internal.AnyKeyPrev
}

func IsComboJustPressed(keys ...int) bool {
	return combo(keys) && IsKeyJustPressed(keys[len(keys)-1])
}
func IsComboHeld(keys ...int) bool {
	return combo(keys) && IsKeyHeld(keys[len(keys)-1])
}

// private ========================================================

func combo(keys []int) bool {
	for _, k := range keys {
		if k < 0 || k >= len(internal.Keys) || !internal.Keys[k] {
			return false
		}
	}

	count := 0
	for _, pressed := range internal.Keys {
		if pressed {
			count++
		}
	}
	return count == len(keys)
}
