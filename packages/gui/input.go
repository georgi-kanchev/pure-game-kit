package gui

import (
	"pure-game-kit/packages/geometry"
	"pure-game-kit/packages/input/mouse"
	"pure-game-kit/packages/input/mouse/button"
	"pure-game-kit/packages/input/mouse/cursor"
	"pure-game-kit/packages/internal"
)

func IsAnyTyping() bool { return typingIn != 0 }
func IsTyping() bool    { return typingIn == widgetCounter }

func IsHovered() bool       { return nowHovered == widgetCounter }
func IsFocused() bool       { return nowFocused == widgetCounter }
func IsJustFocused() bool   { return nowFocused == widgetCounter && lastFocused != widgetCounter }
func IsJustUnfocused() bool { return lastFocused == widgetCounter && nowFocused != widgetCounter }

func IsClicked() bool     { return clickedWidget == widgetCounter }
func IsJustClicked() bool { return justClickedWidget == widgetCounter }

func IsJustDragged() bool { return IsFocused() && mouse.IsButtonJustPressed(button.Left) }
func IsJustDropped() bool {
	return lastClickedWidget == widgetCounter && !IsClicked() && mouse.IsButtonJustReleased(button.Left)
}
func IsJustDroppedUpon() bool {
	return drag != (Area{}) && IsHovered() && mouse.IsButtonJustReleased(button.Left)
}
func Drag() Area {
	if IsJustDragged() {
		drag = widgetArea
	}
	if IsClicked() {
		var mx, my = mouse.CursorDelta()
		drag.X, drag.Y = drag.X+mx/Scale, drag.Y+my/Scale
	}
	return drag
}

// private ========================================================

var skipInput bool    // used for internal calls to the widget input functions only for drawing (no input)
var widgetCounter int // resets every frame, each widget increases it, used for id, checked against the below ids

var nowHovered, lastHovered, nowFocused, lastFocused int
var lastClickedWidget, clickedWidget, justClickedWidget int
var scrollDraggedWidget, scrollHoveredWidget, lastScrollHoveredWidget int
var lastUpdateOnFrame, lastScrollFrame uint64
var widgetArea, drag Area
var droppedLastFrame bool

var lastTypingIn, typingIn, inputIndexCursor, inputIndexSelection int
var inputCursorTimer, ax, bx, inputScroll float32

func handleInput(area, mask Area, roundness float32) {
	if skipInput {
		return
	}

	if internal.Frame != lastUpdateOnFrame { // frame reset, runs exactly once on the first widget of a new frame
		lastUpdateOnFrame, lastFocused, lastTypingIn = internal.Frame, nowFocused, typingIn
		inputCursorTimer += internal.FrameDelta

		if nowHovered == lastHovered {
			nowFocused = nowHovered // widget won input last frame AND won input this frame
		} else {
			nowFocused = 0 // focus broken or shifting
		}

		lastClickedWidget, justClickedWidget = clickedWidget, 0
		if mouse.IsButtonJustPressed(button.Left) {
			clickedWidget = nowFocused //  lock the active widget to whatever is currently hovered
		} else if mouse.IsButtonJustReleased(button.Left) {
			if clickedWidget != 0 && clickedWidget == nowFocused { // same widget we started clicking on?
				justClickedWidget = clickedWidget
			}
			clickedWidget = 0 // clear the lock
		} else if !mouse.IsButtonPressed(button.Left) {
			clickedWidget = 0 // if the button not held, ensure nothing is active
		}

		if droppedLastFrame {
			drag = Area{}
			droppedLastFrame = false
		}
		if mouse.IsButtonJustReleased(button.Left) && drag != (Area{}) {
			droppedLastFrame = true
		}

		lastHovered = nowHovered
		view.Zoom, nowHovered, widgetCounter = Scale, 0, 0
		mouse.SetCursor(cursor.Arrow)
	}

	widgetCounter, widgetArea = widgetCounter+1, area

	var mx, my = view.MousePosition()
	var shape = geometry.NewRoundedRectangle(area.X, area.Y, area.Width, area.Height, 0, roundness)
	var maskHor = mx >= mask.X-mask.Width/2 && mx <= mask.X+mask.Width/2
	var maskVer = my >= mask.Y-mask.Height/2 && my <= mask.Y+mask.Height/2
	var maskCheck = mask == (Area{}) || (mask != (Area{}) && maskHor && maskVer)
	if shape.ContainsPoint(mx, my) && maskCheck {
		nowHovered = widgetCounter // top-most logic: later widgets naturally overwrite earlier widgets
	}
}
