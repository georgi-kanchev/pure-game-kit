package keyboard

import (
	"pure-kit/engine/execution/condition"
	"pure-kit/engine/internal"

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
	return condition.TrueOnce(len(internal.Keys) > 0, ";;keyboard-any-pressed")
}
