package mouse

import (
	"pure-kit/engine/execution/condition"
	"pure-kit/engine/internal"

	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	ButtonLeft = iota
	ButtonRight
	ButtonMiddle
	ButtonExtra1
	ButtonExtra2
)

const (
	CursorDefault = iota
	CursorArrow
	CursorInput
	CursorCrosshair
	CursorHand
	CursorResize1
	CursorResize2
	CursorResize3
	CursorResize4
	CursorMove
	CursorNotAllowed
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

func Scroll() int {
	return int(rl.GetMouseWheelMoveV().Y)
}
func ButtonsPressed() []int {
	return internal.Buttons
}

func IsButtonPressed(button int) bool {
	return rl.IsMouseButtonDown(rl.MouseButton(button))
}
func IsButtonPressedOnce(button int) bool {
	return rl.IsMouseButtonPressed(rl.MouseButton(button))
}
func IsButtonReleasedOnce(button int) bool {
	return rl.IsMouseButtonReleased(rl.MouseButton(button))
}

func IsAnyButtonPressed() bool {
	return len(internal.Buttons) > 0
}
func IsAnyButtonPressedOnce() bool {
	return condition.TrueOnce(len(internal.Buttons) > 0, ";;mouse-any-pressed")
}
