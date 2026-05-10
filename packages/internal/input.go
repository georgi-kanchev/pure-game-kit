package internal

import (
	"pure-game-kit/packages/utility/number"

	rl "github.com/gen2brain/raylib-go/raylib"
)

var Cursor int
var Input = ""
var MouseX, MouseY, MouseDeltaX, MouseDeltaY, Scroll, SmoothScroll float32

var Keys, KeysPrev [350]bool
var KeyCount int

var Buttons, ButtonsPrev [7]bool
var AnyButton, AnyButtonPrev, AnyKey, AnyKeyPrev bool

func AccumulateInput() {
	for i := range 7 {
		if rl.IsMouseButtonDown(rl.MouseButton(i)) {
			if !buttons[i] {
				buttons[i] = true
				activeButtons = append(activeButtons, i)
			}
			anyButton = true
		}
	}

	for {
		var key = int(rl.GetKeyPressed())
		if key <= 0 {
			break
		}
		if key < len(keys) && !keys[key] {
			keys[key] = true
			activeKeys = append(activeKeys, key)
			anyKey = true
		}
	}

	for {
		var char = rl.GetCharPressed()
		if char == 0 {
			break
		}
		input += string(char)
	}

	scroll += rl.GetMouseWheelMoveV().Y
	var pos = rl.GetMousePosition()
	mouseX, mouseY = pos.X, pos.Y

	if prevCursor != Cursor {
		rl.SetMouseCursor(int32(Cursor))
		prevCursor = Cursor
	}
}

func SyncAccumulatedInput() {
	Input, input = input, ""

	AnyKeyPrev, AnyKey = AnyKey, anyKey
	KeysPrev, Keys = Keys, keys
	KeyCount = len(activeKeys) // Set the count for the keyboard package

	for _, k := range activeKeys {
		keys[k] = false
	}
	activeKeys = activeKeys[:0]
	anyKey = false

	//=================================================================

	prevMouseX, prevMouseY = MouseX, MouseY
	MouseX, MouseY = mouseX, mouseY
	MouseDeltaX, MouseDeltaY = MouseX-prevMouseX, MouseY-prevMouseY

	AnyButtonPrev, AnyButton = AnyButton, anyButton
	ButtonsPrev, Buttons = Buttons, buttons

	for _, b := range activeButtons {
		buttons[b] = false
	}
	activeButtons = activeButtons[:0]
	anyButton = false

	//=================================================================

	Scroll, scroll = scroll, 0

	const scrollAccel, scrollDecay = 600.0, 8.0
	SmoothScroll += Scroll * scrollAccel * TickDelta
	SmoothScroll *= number.Exponential(-scrollDecay * TickDelta)

	if SmoothScroll != 0 && number.IsWithin(SmoothScroll, 0, 0.0001) {
		SmoothScroll = 0
	}
	if AnyButton || AnyKey {
		SmoothScroll = 0
	}

	if !WindowFocused {
		Keys, Buttons = [350]bool{}, [7]bool{}
		KeysPrev, ButtonsPrev = [350]bool{}, [7]bool{}
		AnyKey, AnyButton, KeyCount = false, false, 0
	}
}

// private ========================================================

var input string
var scroll float32
var mouseX, mouseY float32
var buttons [7]bool
var keys [350]bool
var anyKey, anyButton bool

var activeKeys []int
var activeButtons []int
var prevMouseX, prevMouseY float32
var prevCursor int
