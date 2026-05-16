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

func (b *Bus) AccumulateInput() {
	for i := range 5 {
		var btn = rl.MouseButton(i)
		if rl.IsMouseButtonPressed(btn) {
			btns[i] = true
			activeBtns = append(activeBtns, i)
		}
	}
	for i := len(activeBtns) - 1; i >= 0; i-- {
		var btn = activeBtns[i]
		if rl.IsMouseButtonReleased(rl.MouseButton(btn)) {
			btns[btn] = false
		}
	}

	scroll += rl.GetMouseWheelMoveV().Y
	var pos = rl.GetMousePosition()
	mouseX, mouseY = pos.X, pos.Y

	// cursor comes from the ticker via the bus (no race)
	var cursor int
	if b != nil {
		cursor = b.Cursor
	}
	if prevCursor != cursor {
		if cursor == -1 {
			rl.HideCursor()
		} else {
			rl.ShowCursor()
			rl.SetMouseCursor(int32(cursor))
			prevCursor = cursor
		}
	}

	//=================================================================

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
		input += string(char)
	}

	windowFocused := rl.IsWindowFocused()
	if !windowFocused {
		btns, keys = [5]bool{}, [350]bool{}
		activeBtns, activeKeys = activeBtns[:0], activeKeys[:0]
	}

	// cleanup released keys/buttons (was in SyncAccumulatedInput — now on the
	// main thread so the active slices are never touched by the ticker)
	for i := len(activeBtns) - 1; i >= 0; i-- {
		if !btns[activeBtns[i]] {
			activeBtns = slices.Delete(activeBtns, i, i+1)
		}
	}
	for i := len(activeKeys) - 1; i >= 0; i-- {
		var key = activeKeys[i]
		if !keys[key] {
			KeyDurs[key] = 0
			activeKeys = slices.Delete(activeKeys, i, i+1)
		}
	}

	accumWindowFocused = windowFocused
}

func (b *Bus) CopyInputToBus() {
	var snap = &b.InputSnap
	snap.MouseX, snap.MouseY = mouseX, mouseY
	snap.Scroll, scroll = scroll, 0
	snap.Input, input = input, ""
	snap.ActiveBtns = append(snap.ActiveBtns[:0], activeBtns...)
	snap.ActiveKeys = append(snap.ActiveKeys[:0], activeKeys...)
	snap.KeyDurs = KeyDurs
	snap.WindowFocused = accumWindowFocused
}

func (b *Bus) SyncAccumulatedInput() {
	var snap = &b.InputSnap

	prevMouseX, prevMouseY = MouseX, MouseY
	MouseX, MouseY = snap.MouseX, snap.MouseY
	MouseDeltaX, MouseDeltaY = MouseX-prevMouseX, MouseY-prevMouseY

	AnyBtnPrev, BtnsPrev = AnyBtn, Btns
	AnyBtn = len(snap.ActiveBtns) > 0

	Btns = [5]bool{}
	for _, btn := range snap.ActiveBtns {
		Btns[btn] = true
	}

	Scroll = snap.Scroll

	const scrollAccel, scrollDecay = 600.0, 8.0
	SmoothScroll += Scroll * scrollAccel * TickDelta
	SmoothScroll *= number.Exponential(-scrollDecay * TickDelta)

	if SmoothScroll != 0 && number.IsWithin(SmoothScroll, 0, 0.0001) {
		SmoothScroll = 0
	}
	if AnyBtn || AnyKey {
		SmoothScroll = 0
	}

	//=================================================================

	Input = snap.Input

	AnyKeyPrev, KeysPrev = AnyKey, Keys
	KeyCount = len(snap.ActiveKeys)
	AnyKey = KeyCount > 0

	Keys = [350]bool{}
	KeyDurs = snap.KeyDurs
	for _, key := range snap.ActiveKeys {
		Keys[key] = true
	}

	//=================================================================

	if !snap.WindowFocused {
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
var prevMouseX, prevMouseY float32
var prevCursor int
var accumWindowFocused bool
