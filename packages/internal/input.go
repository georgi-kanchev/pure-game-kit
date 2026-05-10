// Input synchronization system:
// Bridging the gap between the high-frequency renderer loop and the fixed-frequency ticker loop.
// 1. UpdateInputFromRenderer: "Capture" all events (presses, releases, movement) into accumulators.
// 2. SyncInputFromTicker: "Harvest" accumulators into public variables at the start of every logic tick.
// This ensures fast user interactions (like sub-tick clicks) are never missed by the game logic.

package internal

import (
	"pure-game-kit/packages/utility/collection"
	"pure-game-kit/packages/utility/number"

	rl "github.com/gen2brain/raylib-go/raylib"
)

var Cursor int
var Input = ""
var MouseX, MouseY, MouseDeltaX, MouseDeltaY, Scroll, SmoothScroll float32
var Keys, KeysPrev, Buttons, ButtonsPrev = []int{}, []int{}, []int{}, []int{}
var AnyButtonJustPressed, AnyButtonJustReleased, AnyKeyJustPressed, AnyKeyJustReleased = false, false, false, false

func UpdateInputFromRenderer() {
	for i := range 7 {
		if rl.IsMouseButtonPressed(rl.MouseButton(i)) {
			accumulatedButtonsPressed[i] = true
		}
		if rl.IsMouseButtonReleased(rl.MouseButton(i)) {
			accumulatedButtonsReleased[i] = true
		}
	}

	checkKeyRange(32, 96)
	checkKeyRange(256, 349)

	for {
		var char = rl.GetCharPressed()
		if char == 0 {
			break
		}
		accumulatedInput += string(char)
	}

	accumulatedScroll += rl.GetMouseWheelMoveV().Y
	var delta = rl.GetMouseDelta()
	accumulatedMouseDeltaX += delta.X
	accumulatedMouseDeltaY += delta.Y
	var pos = rl.GetMousePosition()
	MouseX, MouseY = pos.X, pos.Y

	if prevCursor != Cursor {
		rl.SetMouseCursor(int32(Cursor))
	}
	prevCursor = Cursor
}
func SyncInputFromTicker() {
	Input, accumulatedInput = accumulatedInput, ""
	Scroll, accumulatedScroll = accumulatedScroll, 0
	MouseDeltaX, accumulatedMouseDeltaX = accumulatedMouseDeltaX, 0
	MouseDeltaY, accumulatedMouseDeltaY = accumulatedMouseDeltaY, 0

	AnyButtonJustPressed, AnyButtonJustReleased = false, false
	ButtonsPrev = append(ButtonsPrev[:0], Buttons...)
	for i, pressed := range accumulatedButtonsPressed {
		if pressed {
			if !collection.Contains(Buttons, i) {
				Buttons = append(Buttons, i)
			}
			accumulatedButtonsPressed[i], AnyButtonJustPressed = false, true
		}
	}
	for i, released := range accumulatedButtonsReleased {
		if released {
			Buttons = collection.Remove(Buttons, i)
			accumulatedButtonsReleased[i], AnyButtonJustReleased = false, true
		}
	}

	AnyKeyJustPressed, AnyKeyJustReleased = false, false
	KeysPrev = append(KeysPrev[:0], Keys...)
	for i, pressed := range accumulatedKeysPressed {
		if pressed {
			if !collection.Contains(Keys, i) {
				Keys = append(Keys, i)
			}
			accumulatedKeysPressed[i], AnyKeyJustPressed = false, true
		}
	}
	for i, released := range accumulatedKeysReleased {
		if released {
			Keys = collection.Remove(Keys, i)
			accumulatedKeysReleased[i], AnyKeyJustReleased = false, true
		}
	}

	if !WindowFocused {
		Keys, Buttons = Keys[:0], Buttons[:0]
	}

	const scrollAccel, scrollDecay = 600.0, 8.0
	SmoothScroll += Scroll * scrollAccel * TickDelta
	SmoothScroll *= number.Exponential(-scrollDecay * TickDelta)
	if SmoothScroll != 0 && number.IsWithin(SmoothScroll, 0, 0.0001) {
		SmoothScroll = 0
	}
	if AnyButtonJustPressed || AnyButtonJustReleased || AnyKeyJustPressed || AnyKeyJustReleased {
		SmoothScroll = 0
	}
}

// private ========================================================

var accumulatedInput string
var accumulatedScroll, accumulatedMouseDeltaX, accumulatedMouseDeltaY float32
var accumulatedButtonsPressed, accumulatedButtonsReleased [7]bool
var accumulatedKeysPressed, accumulatedKeysReleased [350]bool
var prevCursor int

func checkKeyRange(from, to int) {
	for i := from; i < to+1; i++ {
		if rl.IsKeyPressed(int32(i)) {
			accumulatedKeysPressed[i] = true
		}
		if rl.IsKeyReleased(int32(i)) {
			accumulatedKeysReleased[i] = true
		}
	}
}
