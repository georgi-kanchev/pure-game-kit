package internal

import (
	"pure-game-kit/packages/utility/number"
	"slices"

	rl "github.com/gen2brain/raylib-go/raylib"
)

var Cursor int
var Input = ""
var MouseX, MouseY, MouseDeltaX, MouseDeltaY, Scroll, SmoothScroll float32

var Keys, KeysPrev [350]bool
var KeyDurs [350]float32
var AnyKey, AnyKeyPrev bool
var KeyCount int

var Btns, BtnsPrev [5]bool
var AnyBtn, AnyBtnPrev bool

func UpdateInput() {
	Input = ""
	Scroll = 0
	updateKeyboard()
	updateMouse()
}

func updateMouse() {
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

	Scroll += rl.GetMouseWheelMoveV().Y
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

	// cleanup released buttons
	for i := len(activeBtns) - 1; i >= 0; i-- {
		if !btns[activeBtns[i]] {
			activeBtns = slices.Delete(activeBtns, i, i+1)
		}
	}

	AnyBtn = len(activeBtns) > 0
	Btns = [5]bool{}
	for _, btn := range activeBtns {
		Btns[btn] = true
	}

	const scrollAccel, scrollDecay = 600.0, 8.0
	SmoothScroll += Scroll * scrollAccel * TickDelta
	SmoothScroll *= number.Exponential(-scrollDecay * TickDelta)
	if SmoothScroll != 0 && number.IsWithin(SmoothScroll, 0, 0.0001) {
		SmoothScroll = 0
	}
	if AnyBtn || AnyKey {
		SmoothScroll = 0
	}

	if !WindowFocused {
		Btns, BtnsPrev = [5]bool{}, [5]bool{}
		AnyBtn, AnyBtnPrev = false, false
		btns = [5]bool{}
		activeBtns = activeBtns[:0]
	}
}

func updateKeyboard() {
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
		Input += string(char)
	}

	// cleanup released keys
	for i := len(activeKeys) - 1; i >= 0; i-- {
		var key = activeKeys[i]
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

// private ========================================================

var btns [5]bool
var keys [350]bool

var activeKeys []int32
var activeBtns []int
var prevCursor int
