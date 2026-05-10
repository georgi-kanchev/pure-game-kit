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

func UpdateInput() {
	AnyButtonJustPressed = false
	AnyButtonJustReleased = false
	ButtonsPrev = append(ButtonsPrev[:0], Buttons...)

	for i := range 7 {
		if rl.IsMouseButtonPressed(rl.MouseButton(i)) {
			Buttons = append(Buttons, i)
			AnyButtonJustPressed = true
		}
		if rl.IsMouseButtonReleased(rl.MouseButton(i)) {
			Buttons = collection.Remove(Buttons, i)
			AnyButtonJustReleased = true
		}
	}

	if prevCursor != Cursor {
		rl.SetMouseCursor(int32(Cursor))
	}
	prevCursor = Cursor

	var delta = rl.GetMouseDelta()
	var pos = rl.GetMousePosition()
	MouseDeltaX, MouseDeltaY = delta.X, delta.Y
	MouseX, MouseY = pos.X, pos.Y
	Scroll = rl.GetMouseWheelMoveV().Y

	const scrollAccel, scrollDecay = 600.0, 8.0
	SmoothScroll += Scroll * scrollAccel * TickDelta
	SmoothScroll *= number.Exponential(-scrollDecay * TickDelta)
	if SmoothScroll != 0 && number.IsWithin(SmoothScroll, 0, 0.0001) {
		SmoothScroll = 0
	}
	if AnyButtonJustPressed || AnyButtonJustReleased || AnyKeyJustPressed || AnyKeyJustReleased {
		SmoothScroll = 0
	}

	//=================================================================

	AnyKeyJustPressed = false
	AnyKeyJustReleased = false
	KeysPrev = append(KeysPrev[:0], Keys...)
	Input = ""

	var char = rl.GetCharPressed()
	for char > 0 {
		Input += string(char)
		char = rl.GetCharPressed()
	}

	checkKeyRange(32, 96)
	checkKeyRange(256, 349)

	if !WindowFocused {
		Keys = Keys[:0]
		Buttons = Buttons[:0]
	}
}

// private ========================================================

func checkKeyRange(from, to int) {
	for i := from; i < to+1; i++ {
		if rl.IsKeyPressed(int32(i)) {
			Keys = append(Keys, i)
			AnyKeyJustPressed = true
		}
		if rl.IsKeyReleased(int32(i)) {
			Keys = collection.Remove(Keys, i)
			AnyKeyJustReleased = true
		}
	}
}
