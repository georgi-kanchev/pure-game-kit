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
var KeyCount int // for the combo

var Btns, BtnsPrev [5]bool
var AnyBtn, AnyBtnPrev bool

// AccumulateAndSyncInput collects raw input from raylib and publishes it to the
// public vars. Called once per frame before game logic. Single-threaded — no
// snapshot or cross-goroutine synchronization needed.
func AccumulateAndSyncInput() {
	// ---- accumulate raw events (was DoInput) ----

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

	scroll += rl.GetMouseWheelMoveV().Y
	var pos = rl.GetMousePosition()
	mouseX, mouseY = pos.X, pos.Y

	if prevCursor != Cursor {
		if Cursor == -1 {
			rl.HideCursor()
		} else {
			rl.ShowCursor()
			rl.SetMouseCursor(int32(Cursor))
			prevCursor = Cursor
		}
	}

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
		accumKeyDurs[key] += FrameDelta
		if rl.IsKeyReleased(key) {
			keys[key] = false
		}
	}

	for {
		var char = rl.GetCharPressed()
		if char == 0 {
			break
		}
		input += string(char)
	}

	windowFocused := rl.IsWindowFocused()
	if !windowFocused {
		btns, keys = [5]bool{}, [350]bool{}
		activeBtns, activeKeys = activeBtns[:0], activeKeys[:0]
	}

	// cleanup released keys/buttons
	for i := len(activeBtns) - 1; i >= 0; i-- {
		if !btns[activeBtns[i]] {
			activeBtns = slices.Delete(activeBtns, i, i+1)
		}
	}
	for i := len(activeKeys) - 1; i >= 0; i-- {
		var key = activeKeys[i]
		if !keys[key] {
			accumKeyDurs[key] = 0
			activeKeys = slices.Delete(activeKeys, i, i+1)
		}
	}

	// ---- publish to public vars (was UnpackForTickInput) ----

	prevMouseX, prevMouseY = MouseX, MouseY
	MouseX, MouseY = mouseX, mouseY
	MouseDeltaX, MouseDeltaY = MouseX-prevMouseX, MouseY-prevMouseY

	AnyBtnPrev, BtnsPrev = AnyBtn, Btns
	AnyBtn = len(activeBtns) > 0
	Btns = [5]bool{}
	for _, btn := range activeBtns {
		Btns[btn] = true
	}

	Scroll, scroll = scroll, 0

	const scrollAccel, scrollDecay = 600.0, 8.0
	SmoothScroll += Scroll * scrollAccel * TickDelta
	SmoothScroll *= number.Exponential(-scrollDecay * TickDelta)
	if SmoothScroll != 0 && number.IsWithin(SmoothScroll, 0, 0.0001) {
		SmoothScroll = 0
	}
	if AnyBtn || AnyKey {
		SmoothScroll = 0
	}

	Input, input = input, ""

	AnyKeyPrev, KeysPrev = AnyKey, Keys
	KeyCount = len(activeKeys)
	AnyKey = KeyCount > 0
	Keys = [350]bool{}
	KeyDurs = accumKeyDurs
	for _, key := range activeKeys {
		Keys[key] = true
	}

	if !windowFocused {
		Btns, Keys = [5]bool{}, [350]bool{}
		BtnsPrev, KeysPrev = [5]bool{}, [350]bool{}
		AnyBtn, AnyKey, KeyCount = false, false, 0
	}
}

// private ========================================================

var input string
var scroll float32
var mouseX, mouseY float32
var btns [5]bool
var keys [350]bool

var activeKeys []int32
var activeBtns []int
var accumKeyDurs [350]float32
var prevMouseX, prevMouseY float32
var prevCursor int
