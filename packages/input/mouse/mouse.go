// Contains checks for whether a certain mouse button is interacted with in various ways.
// Also provides the currently pressed buttons, scrolling and OS cursor customization.
// Meant to be checked every frame instead of subscribtion-based events/callbacks.
package mouse

import (
	"pure-game-kit/packages/internal"

	rl "github.com/gen2brain/raylib-go/raylib"
)

//=================================================================

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

func IsHoveringWindow() bool {
	return internal.WindowHovered
}

func CursorDelta() (x, y float32) {
	return internal.MouseDeltaX, internal.MouseDeltaY
}

func Scroll() float32 {
	return internal.Scroll
}
func ScrollSmooth() float32 {
	return internal.SmoothScroll
}

func IsButtonPressed(button int) bool {
	if button < 0 || button >= len(internal.Buttons) {
		return false
	}
	return internal.Buttons[button]
}
func IsButtonJustPressed(button int) bool {
	if button < 0 || button >= len(internal.Buttons) {
		return false
	}
	return internal.Buttons[button] && !internal.ButtonsPrev[button]
}
func IsButtonJustReleased(button int) bool {
	if button < 0 || button >= len(internal.Buttons) {
		return false
	}
	return !internal.Buttons[button] && internal.Buttons[button]
}

func IsAnyButtonPressed() bool {
	return internal.AnyButton
}
func IsAnyButtonJustPressed() bool {
	return internal.AnyButton && !internal.AnyButtonPrev
}
func IsAnyButtonJustReleased() bool {
	return !internal.AnyButton && internal.AnyButtonPrev
}
