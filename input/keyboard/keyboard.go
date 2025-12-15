/*
Contains checks for whether a certain keyboard key is interacted with in various ways.
Also provides text input and the currently pressed keys. Meant to be checked every frame
instead of subscribtion-based events/callbacks.
*/
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
	return collection.Contains(internal.Keys, key)
}
func IsKeyHeld(key int) bool {
	return rl.IsKeyPressedRepeat(int32(key))
}

func IsKeyJustPressed(key int) bool {
	return collection.Contains(internal.Keys, key) && !collection.Contains(internal.KeysPrev, key)
}
func IsKeyJustReleased(key int) bool {
	return !collection.Contains(internal.Keys, key) && collection.Contains(internal.KeysPrev, key)
}

func IsAnyKeyPressed() bool {
	return len(internal.Keys) > 0
}
func IsAnyKeyJustPressed() bool {
	return internal.AnyKeyJustPressed
}
func IsAnyKeyJustReleased() bool {
	return internal.AnyKeyJustReleased
}

func IsComboJustPressed(keys ...int) bool {
	return combo(keys) && IsKeyJustPressed(keys[len(keys)-1])
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
