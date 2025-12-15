/*
Contains checks for whether a certain mouse button is interacted with in various ways.
Also provides the currently pressed buttons, scrolling and OS cursor customization.
Meant to be checked every frame instead of subscribtion-based events/callbacks.
*/
package mouse

import (
	"pure-game-kit/internal"
	"pure-game-kit/utility/collection"

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

func CursorDelta() (x, y float32) {
	return internal.MouseDeltaX, internal.MouseDeltaY
}

func Scroll() float32 {
	return rl.GetMouseWheelMoveV().Y
}
func ScrollSmooth() float32 {
	return internal.SmoothScroll
}

func ButtonsPressed() []int {
	return internal.Buttons
}

func IsButtonPressed(button int) bool {
	return collection.Contains(internal.Buttons, button)
}
func IsButtonJustPressed(button int) bool {
	return collection.Contains(internal.Buttons, button) && !collection.Contains(internal.ButtonsPrev, button)
}
func IsButtonJustReleased(button int) bool {
	return !collection.Contains(internal.Buttons, button) && collection.Contains(internal.ButtonsPrev, button)
}

func IsAnyButtonPressed() bool {
	return len(internal.Buttons) > 0
}
func IsAnyButtonJustPressed() bool {
	return internal.AnyButtonJustPressed
}
func IsAnyButtonJustReleased() bool {
	return internal.AnyButtonJustReleased
}
