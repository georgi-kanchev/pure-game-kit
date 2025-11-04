package mouse

import (
	"pure-game-kit/internal"

	rl "github.com/gen2brain/raylib-go/raylib"
)

//=================================================================
// setters

func SetCursorVisibility(visible bool) {
	if visible {
		rl.ShowCursor()
	} else {
		rl.HideCursor()
	}
}
func SetCursor(cursor int) {
	internal.Cursor = cursor
}

//=================================================================
// getters

func CursorDelta() (x, y float32) {
	var delta = rl.GetMouseDelta()
	return delta.X, delta.Y
}
func Scroll() int {
	return int(rl.GetMouseWheelMoveV().Y)
}
func ButtonsPressed() []int {
	return internal.Buttons
}

func IsButtonPressed(button int) bool {
	return rl.IsMouseButtonDown(rl.MouseButton(button))
}
func IsButtonJustPressed(button int) bool {
	return rl.IsMouseButtonPressed(rl.MouseButton(button))
}
func IsButtonJustReleased(button int) bool {
	return rl.IsMouseButtonReleased(rl.MouseButton(button))
}

func IsAnyButtonPressed() bool {
	return len(internal.Buttons) > 0
}
func IsAnyButtonJustPressed() bool {
	return internal.AnyButtonPressedOnce
}
func IsAnyButtonJustReleased() bool {
	return internal.AnyButtonReleasedOnce
}
