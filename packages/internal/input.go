package internal

import (
	"pure-game-kit/packages/utility/number"
	"slices"

	rl "github.com/gen2brain/raylib-go/raylib"
)

var Cursor int
var Input []rune = make([]rune, 0, 8)
var MouseX, MouseY, MouseDeltaX, MouseDeltaY, ScrollX, SmoothScrollX, ScrollY, SmoothScrollY float32

var Keys, KeysPrev [350]bool
var KeyDurs [350]float32
var AnyKey, AnyKeyPrev bool
var KeyCount int

var Btns, BtnsPrev [5]bool
var AnyBtn, AnyBtnPrev bool

func CacheInput() {
	ScrollX, ScrollY, Input = 0, 0, Input[:0]
	cacheKeyboard()
	cacheMouse()
}

// private ========================================================

var btns [5]bool
var keys [350]bool

var activeKeys []int32 = make([]int32, 0, 8)
var activeBtns []int = make([]int, 0, 8)
var prevCursor int

func cacheMouse() {
	AnyBtnPrev, BtnsPrev = AnyBtn, Btns

	for i := range 5 {
		var btn = rl.MouseButton(i)
		if rl.IsMouseButtonPressed(btn) {
			btns[i] = true
			activeBtns = append(activeBtns, i)
		}
	}
	for i := len(activeBtns) - 1; i >= 0; i-- {
		if rl.IsMouseButtonReleased(rl.MouseButton(activeBtns[i])) {
			btns[activeBtns[i]] = false
		}
	}

	var wheel = rl.GetMouseWheelMoveV()
	ScrollX, ScrollY = ScrollX+wheel.X, ScrollY+wheel.Y
	var pos = rl.GetMousePosition()
	MouseDeltaX, MouseDeltaY = pos.X-MouseX, pos.Y-MouseY
	MouseX, MouseY = pos.X, pos.Y

	if prevCursor != Cursor {
		if Cursor == -1 {
			rl.HideCursor()
		} else {
			rl.ShowCursor()
			rl.SetMouseCursor(int32(Cursor))
			prevCursor = Cursor
		}
	}

	for i := len(activeBtns) - 1; i >= 0; i-- {
		if !btns[activeBtns[i]] { // cleanup released buttons
			activeBtns = slices.Delete(activeBtns, i, i+1)
		}
	}

	AnyBtn = len(activeBtns) > 0
	Btns = [5]bool{}
	for _, btn := range activeBtns {
		Btns[btn] = true
	}

	const scrollAccel, scrollDecay = 600.0, 8.0
	SmoothScrollX += ScrollX * scrollAccel * FrameDelta
	SmoothScrollY += ScrollY * scrollAccel * FrameDelta
	SmoothScrollX *= number.Exponential(-scrollDecay * FrameDelta)
	SmoothScrollY *= number.Exponential(-scrollDecay * FrameDelta)
	if SmoothScrollX != 0 && number.IsWithin(SmoothScrollX, 0, 0.0001) {
		SmoothScrollX = 0
	}
	if SmoothScrollY != 0 && number.IsWithin(SmoothScrollY, 0, 0.0001) {
		SmoothScrollY = 0
	}
	if AnyBtn {
		SmoothScrollX = 0
		SmoothScrollY = 0
	}

	if !WindowFocused {
		Btns, BtnsPrev = [5]bool{}, [5]bool{}
		AnyBtn, AnyBtnPrev = false, false
		btns = [5]bool{}
		activeBtns = activeBtns[:0]
	}
}
func cacheKeyboard() {
	AnyKeyPrev, KeysPrev = AnyKey, Keys

	for {
		var key = rl.GetKeyPressed()
		if key <= 0 || key >= 350 {
			break
		}
		keys[key] = true
		activeKeys = append(activeKeys, key)
	}
	for i := len(activeKeys) - 1; i >= 0; i-- {
		var key = activeKeys[i]
		KeyDurs[key] += FrameDelta
		if rl.IsKeyReleased(key) {
			keys[key] = false
		}
	}

	for {
		var char = rl.GetCharPressed()
		if char == 0 {
			break
		}
		Input = append(Input, char)
	}

	for i := len(activeKeys) - 1; i >= 0; i-- {
		var key = activeKeys[i] // cleanup released keys
		if !keys[key] {
			KeyDurs[key] = 0
			activeKeys = slices.Delete(activeKeys, i, i+1)
		}
	}

	KeyCount = len(activeKeys)
	AnyKey = KeyCount > 0
	Keys = [350]bool{}
	for _, key := range activeKeys {
		Keys[key] = true
	}

	if !WindowFocused {
		Keys, KeysPrev = [350]bool{}, [350]bool{}
		AnyKey, AnyKeyPrev, KeyCount = false, false, 0
		keys = [350]bool{}
		activeKeys = activeKeys[:0]
	}
}
