package mouse

import i "pure-game-kit/packages/internal"

func SetCursor(cursor int) { i.Cursor = cursor }

func CursorDelta() (x, y float32) { return i.MouseDeltaX, i.MouseDeltaY }

func ScrollY() float32       { return i.ScrollY }
func ScrollSmoothY() float32 { return i.SmoothScrollY }
func ScrollX() float32       { return i.ScrollX }       // laptop touchpad has it
func ScrollSmoothX() float32 { return i.SmoothScrollX } // laptop touchpad has it

func IsButtonPressed(button int) bool      { return i.Btns[button] }
func IsButtonJustPressed(button int) bool  { return i.Btns[button] && !i.BtnsPrev[button] }
func IsButtonJustReleased(button int) bool { return !i.Btns[button] && i.BtnsPrev[button] }

func IsAnyButtonPressed() bool      { return i.AnyBtn }
func IsAnyButtonJustPressed() bool  { return i.AnyBtn && !i.AnyBtnPrev }
func IsAnyButtonJustReleased() bool { return !i.AnyBtn && i.AnyBtnPrev }
