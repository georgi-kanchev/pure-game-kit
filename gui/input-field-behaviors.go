package gui

import (
	"pure-game-kit/execution/condition"
	"pure-game-kit/gui/field"
	"pure-game-kit/input/keyboard"
	"pure-game-kit/input/keyboard/key"
	"pure-game-kit/input/mouse"
	btn "pure-game-kit/input/mouse/button"
	"pure-game-kit/internal"
	"pure-game-kit/utility/is"
	"pure-game-kit/utility/number"
	txt "pure-game-kit/utility/text"
	"pure-game-kit/utility/time"
)

var clickTime float32
var clickIndex, clickCycle int

func tryMoveCursor(text string, margin float32) {
	var w = typingIn
	var ctrl = keyboard.IsKeyPressed(key.LeftControl) || keyboard.IsKeyPressed(key.RightControl)
	var home = keyboard.IsKeyJustPressed(key.UpArrow) || keyboard.IsKeyJustPressed(key.Home)
	var end = keyboard.IsKeyJustPressed(key.DownArrow) || keyboard.IsKeyJustPressed(key.End)
	var length = txt.Length(text)

	if keyboard.IsKeyJustPressed(key.LeftArrow) || keyboard.IsKeyHeld(key.LeftArrow) {
		var max = number.Biggest(indexCursor-1, 0)
		cursorTime = 0
		indexCursor = condition.If(ctrl, wordIndex(text, true, indexCursor), max)
		trySelect()
	}
	if keyboard.IsKeyJustPressed(key.RightArrow) || keyboard.IsKeyHeld(key.RightArrow) {
		var min = number.Smallest(length, indexCursor+1)
		cursorTime = 0
		indexCursor = condition.If(ctrl, wordIndex(text, false, indexCursor), min)
		trySelect()
	}
	if home || end {
		cursorTime = 0
		indexCursor = condition.If(end, length, 0)
		trySelect()
	}

	if mouse.IsButtonPressed(btn.Left) {
		cursorTime = 0

		if length == 0 {
			symbolXs = []float32{}
			indexCursor = 0
		} else {
			indexCursor = closestIndexToMouse(w.root.cam)

			if mouse.IsButtonJustPressed(btn.Left) {
				calculateXs(w) // calculate once and update indexes to not drop performance
				indexCursor = closestIndexToMouse(w.root.cam)
				indexSelect = indexCursor
			}
		}
	}

	var cx = cursorX(margin, w)
	var left, right = w.X + margin, w.X + w.Width - margin
	if cx < left && indexCursor >= 0 {
		scrollX -= left - cx
		setupText(margin, w, true)
	}
	if cx > right && indexCursor <= length {
		scrollX += cx - right
		setupText(margin, w, true)
	}
}
func trySelect() {
	if keyboard.IsKeyPressed(key.LeftShift) || keyboard.IsKeyPressed(key.RightShift) {
		return
	}

	var a, b = indexSelect, indexCursor
	indexSelect = indexCursor

	if indexCursor != indexSelect {
		indexCursor = condition.If(a < b, a, b)
		indexSelect = indexCursor
	}
}
func tryRemove(text string, margin float32) string {
	var w = typingIn
	var left, right = w.X + margin, w.X + w.Width - margin
	var ctrl = keyboard.IsKeyPressed(key.LeftControl) || keyboard.IsKeyPressed(key.RightControl)
	var remove = func(back, front int) {
		cursorTime = 0

		if back > 0 && indexCursor == 0 {
			return
		}
		if front > 0 && indexCursor == txt.Length(text) {
			return
		}
		text = txt.Part(text, 0, indexCursor-back) + txt.Part(text, indexCursor+front, txt.Length(text))
		setText(w, text)
		indexCursor -= back
		indexSelect = indexCursor
		calculateXs(w)

		var owner = w.root.Containers[w.OwnerId]
		sound.AssetId = defaultValue(w.root.themedField(field.InputFieldSoundErase, owner, w), "~erase")
		sound.Volume = w.root.Volume
		sound.Play()
	}

	if keyboard.IsKeyJustPressed(key.Backspace) || keyboard.IsKeyJustPressed(key.Delete) || simulateRemove {
		if indexSelect < indexCursor {
			remove(indexCursor-indexSelect, 0)
			return text
		} else if indexCursor < indexSelect {
			remove(0, indexSelect-indexCursor)
			return text
		}
	}

	if keyboard.IsKeyJustPressed(key.Backspace) || keyboard.IsKeyHeld(key.Backspace) {
		remove(condition.If(ctrl, indexCursor-wordIndex(text, true, indexCursor), 1), 0)

		// scrolls left when empty space appears on the right (if possible)
		var textWidth, _ = w.textBox.TextMeasure(w.textBox.Text)
		var textRight = (left - scrollX) + textWidth
		if indexCursor > 0 && textRight < right {
			scrollX -= right - textRight
			scrollX = condition.If(textWidth < right-left, 0, scrollX)
			setupText(margin, w, true)
		}
	}
	if keyboard.IsKeyJustPressed(key.Delete) || keyboard.IsKeyHeld(key.Delete) {
		remove(0, condition.If(ctrl, wordIndex(text, false, indexCursor)-indexCursor, 1))
	}
	return text
}
func tryInput(text string, margin float32) string {
	var w = typingIn
	var input = keyboard.Input()
	if is.AnyOf(input, "", "{", "}") {
		return text
	}

	if indexCursor != indexSelect { // text is selected, we should remove it and then type
		simulateRemove = true
		text = tryRemove(text, margin)
		simulateRemove = false
	}

	if txt.Length(text) == 0 {
		text = input
		setText(w, text)
		setupText(margin, w, true) // text is not setuped cuz it was empty "" (skipped)
	} else {
		text = txt.Insert(text, input, indexCursor)
	}

	var owner = w.root.Containers[w.OwnerId]
	sound.AssetId = defaultValue(w.root.themedField(field.InputFieldSoundType, owner, w), "~write")
	sound.Volume = w.root.Volume
	sound.Play()

	setText(w, text)
	indexCursor += txt.Length(input)
	indexSelect = indexCursor
	cursorTime = 0
	calculateXs(w)
	return text
}
func tryFocusNextField() {
	if !keyboard.IsKeyJustPressed(key.Tab) || frame == int(time.FrameCount()) {
		return
	}

	var self = typingIn
	var owner = self.root.Containers[self.OwnerId]
	var allInputFields = []*widget{}
	var myIndex = 0
	for _, wId := range owner.Widgets {
		var w = self.root.Widgets[wId]

		if w.Class == "inputField" && !w.IsCulled && !w.isHidden(owner) {
			allInputFields = append(allInputFields, w)
		}
		if w == self {
			myIndex = len(allInputFields) - 1
		}
	}
	var total = len(allInputFields)
	if total == 1 {
		return // i'm the only input field, do nothing
	}
	cursorTime = 0
	scrollX = 0
	typingIn = allInputFields[(myIndex+1)%total]
	var text = txt.Remove(self.root.themedField(field.Text, owner, typingIn), "\n")
	indexCursor = len(text)
	indexSelect = indexCursor
	frame = int(time.FrameCount()) // only once per frame

	var margin = parseNum(self.root.themedField(field.InputFieldMargin, owner, typingIn), 10)
	setupText(margin, typingIn, true)
	if text == "" { // empty text is skipped in setupText so Xs should affect that
		typingIn.textBox.Text = ""
	}
	calculateXs(typingIn)
}
func tryCycleSelection(text string) (success bool) {
	var left = mouse.IsButtonJustPressed(btn.Left)
	var focused = typingIn.isFocused()
	var timeToCycle = internal.Runtime-clickTime < 0.5
	if !left || !focused {
		return clickCycle > 0
	}

	var index = closestIndexToMouse(typingIn.root.cam)
	var sameClickIndex = clickIndex == index
	clickIndex = index
	clickTime = internal.Runtime

	if !timeToCycle || !sameClickIndex {
		clickCycle = 0
		return false
	}

	clickCycle = (clickCycle + 1) % 3

	switch clickCycle {
	case 0:
		indexSelect, indexCursor = index, index
	case 1:
		indexSelect, indexCursor = wordIndex(text, true, index), wordIndex(text, false, index)
	case 2:
		indexSelect, indexCursor = 0, txt.Length(text)
	}
	return true
}
