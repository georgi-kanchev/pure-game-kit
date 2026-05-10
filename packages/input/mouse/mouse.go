package mouse

import i "pure-game-kit/packages/internal"

func SetCursor(cursor int) { i.Cursor = cursor }

func CursorDelta() (x, y float32) { return i.MouseDeltaX, i.MouseDeltaY }

func Scroll() float32       { return i.Scroll }
func ScrollSmooth() float32 { return i.SmoothScroll }

func IsButtonPressed(button int) bool      { return i.Btns[button] }
func IsButtonJustPressed(button int) bool  { return i.Btns[button] && !i.BtnsPrev[button] }
func IsButtonJustReleased(button int) bool { return !i.Btns[button] && i.BtnsPrev[button] }

func IsAnyButtonPressed() bool      { return i.AnyBtn }
func IsAnyButtonJustPressed() bool  { return i.AnyBtn && !i.AnyBtnPrev }
func IsAnyButtonJustReleased() bool { return !i.AnyBtn && i.AnyBtnPrev }
