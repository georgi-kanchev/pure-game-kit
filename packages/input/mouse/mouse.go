package mouse

import "pure-game-kit/packages/internal"

func SetCursor(cursor int) { internal.Cursor = cursor }

func CursorDelta() (x, y float32) { return internal.MouseDeltaX, internal.MouseDeltaY }

func Scroll() float32       { return internal.Scroll }
func ScrollSmooth() float32 { return internal.SmoothScroll }

func IsButtonPressed(button int) bool {
	return button >= 0 && button < len(internal.Buttons) && internal.Buttons[button]
}

func IsButtonJustPressed(button int) bool {
	return IsButtonPressed(button) && !internal.ButtonsPrev[button]
}

func IsButtonJustReleased(button int) bool {
	return button >= 0 && button < len(internal.Buttons) && !internal.Buttons[button] && internal.ButtonsPrev[button]
}

func IsAnyButtonPressed() bool      { return internal.AnyButton }
func IsAnyButtonJustPressed() bool  { return internal.AnyButton && !internal.AnyButtonPrev }
func IsAnyButtonJustReleased() bool { return !internal.AnyButton && internal.AnyButtonPrev }
