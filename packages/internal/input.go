// Input synchronization system:
// Bridging the gap between the high-frequency renderer loop and the fixed-frequency ticker loop.
// 1. UpdateInput (Renderer): "Capture" all events (presses, releases, movement) into accumulators.
// 2. SyncInput (Ticker): "Harvest" accumulators into public variables at the start of every logic tick.
// This ensures fast user interactions (like sub-tick clicks) are never missed by the game logic.
package internal

import (
	"pure-game-kit/packages/utility/number"

	rl "github.com/gen2brain/raylib-go/raylib"
)

var Cursor int
var Input = ""
var MouseX, MouseY, MouseDeltaX, MouseDeltaY, Scroll, SmoothScroll float32

var Keys, KeysPrev [350]bool
var Buttons, ButtonsPrev [7]bool
var AnyButton, AnyButtonPrev, AnyKey, AnyKeyPrev bool

func AccumulateInput() {
	for i := range 7 {
		if rl.IsMouseButtonPressed(rl.MouseButton(i)) {
			buttons[i] = true
			anyButton = true
		}
	}

	for {
		var key = int(rl.GetKeyPressed())
		if key == 0 || key >= len(Keys) {
			break
		}
		keys[key] = true
		anyKey = true
	}

	for {
		var char = rl.GetCharPressed()
		if char == 0 {
			break
		}
		input += string(char)
	}

	scroll += rl.GetMouseWheelMoveV().Y
}
func SyncAccumulatedInput() {
	Input, input = input, ""

	AnyKey, AnyKeyPrev = anyKey, anyKeyPrev
	KeysPrev = Keys // instant copy

	keys = [350]bool{}
	anyKeyPrev = anyKey
	anyKey = false

	//=================================================================

	prevMouseX, prevMouseY = MouseX, MouseY
	var pos = rl.GetMousePosition()
	MouseX, MouseY = pos.X, pos.Y
	MouseDeltaX = MouseX - prevMouseX
	MouseDeltaY = MouseY - prevMouseY

	AnyButton, AnyButtonPrev = anyButton, anyButtonPrev
	ButtonsPrev = Buttons // instant copy

	buttons = [7]bool{}
	anyButtonPrev = anyButton
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

	//=================================================================

	if prevCursor != Cursor {
		rl.SetMouseCursor(int32(Cursor))
		prevCursor = Cursor
	}

	//=================================================================

	if !WindowFocused {
		Keys, Buttons = [350]bool{}, [7]bool{}
		KeysPrev, ButtonsPrev = [350]bool{}, [7]bool{}
		AnyKey, AnyButton = false, false
	}
}

// private ========================================================

var prevMouseX, prevMouseY float32
var prevCursor int

var input string
var scroll float32
var buttons [7]bool
var keys [350]bool
var anyKey, anyKeyPrev, anyButton, anyButtonPrev bool
