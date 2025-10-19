package keyboard

import (
	"pure-game-kit/internal"
	"pure-game-kit/utility/collection"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func Input() string {
	return internal.Input
}
func KeysPressed() []int {
	return internal.Keys
}

func IsKeyPressed(key int) bool {
	return rl.IsKeyDown(int32(key))
}
func IsKeyHeld(key int) bool {
	return rl.IsKeyPressedRepeat(int32(key))
}

func IsKeyPressedOnce(key int) bool {
	return rl.IsKeyPressed(int32(key))
}
func IsKeyReleasedOnce(key int) bool {
	return rl.IsKeyReleased(int32(key))
}

func IsAnyKeyPressed() bool {
	return len(internal.Keys) > 0
}
func IsAnyKeyPressedOnce() bool {
	for _, k := range internal.Keys {
		if !collection.Contains(internal.KeysPrev, k) {
			return true
		}
	}
	return false
}
func IsAnyKeyReleasedOnce() bool {
	for _, k := range internal.KeysPrev {
		if !collection.Contains(internal.Keys, k) {
			return true
		}
	}
	return false
}

func IsComboPressedOnce(keys ...int) bool {
	return combo(keys) && IsKeyPressedOnce(keys[len(keys)-1])
}
func IsComboHeld(keys ...int) bool {
	return combo(keys) && IsKeyHeld(keys[len(keys)-1])
}

//=================================================================
// private

func combo(keys []int) bool {
	if len(internal.Keys) != len(keys) {
		return false
	}

	for i := range internal.Keys {
		if internal.Keys[i] != keys[i] {
			return false
		}
	}
	return true
}
