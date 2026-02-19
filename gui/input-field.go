package gui

import (
	"pure-game-kit/execution/condition"
	"pure-game-kit/graphics"
	"pure-game-kit/gui/field"
	"pure-game-kit/input/keyboard"
	"pure-game-kit/input/keyboard/key"
	"pure-game-kit/input/mouse"
	btn "pure-game-kit/input/mouse/button"
	"pure-game-kit/input/mouse/cursor"
	"pure-game-kit/internal"
	"pure-game-kit/utility/color/palette"
	"pure-game-kit/utility/number"
	txt "pure-game-kit/utility/text"
	"pure-game-kit/utility/time"
	"pure-game-kit/window"
)

func InputField(id string, properties ...string) string {
	return newWidget("inputField", id, properties...)
}

//=================================================================

func (g *GUI) InputFieldTyping() (inputFieldId string) {
	if typingIn != nil {
		return typingIn.Id
	}
	return ""
}
func (g *GUI) InputFieldStopTyping() {
	typingIn = nil
	scrollX = 0
}

//=================================================================
// private

const cursorWidth float32 = 2

func setupText(margin float32, root *root, widget *widget, skipEmpty bool) {
	setupVisualsText(root, widget, skipEmpty)
	widget.textBox.AlignmentX, widget.textBox.AlignmentY = 0, 0
	widget.textBox.Width = 99999
	widget.textBox.Height = widget.Height - margin
	widget.textBox.LineHeight = widget.Height - margin
	var scroll = condition.If(typingIn == widget, scrollX, 0)
	widget.textBox.X = widget.textBox.X + margin - scroll
	widget.textBox.Y += margin / 2
	widget.textBox.Text = txt.Remove(widget.textBox.Text, "\n")
}
func inputField(cam *graphics.Camera, root *root, widget *widget) {
	if widget.textBox == nil {
		widget.textBox = &graphics.TextBox{}
	}

	var owner = root.Containers[widget.OwnerId]
	var margin = parseNum(root.themedField(field.InputFieldMargin, owner, widget), 10)

	if keyboard.IsComboJustPressed(key.LeftControl, key.A) ||
		keyboard.IsComboJustPressed(key.RightControl, key.A) {
		indexSelect = 0
		indexCursor = len(symbolXs) - 1
	}

	setupText(margin, root, widget, true)
	setupVisualsTextured(root, widget)

	if widget.isFocused(root, cam) {
		mouse.SetCursor(cursor.Input)
	}

	var anyInput = mouse.IsAnyButtonJustPressed() || mouse.Scroll() != 0
	var focused = widget.isFocused(root, cam)
	var meTyping = typingIn == widget // each input field should disable its own typing
	var text = txt.Remove(root.themedField(field.Text, owner, widget), "\n")

	if meTyping && ((anyInput && !focused) || !window.IsHovered() || keyboard.IsKeyJustPressed(key.Escape)) {
		typingIn = nil
		scrollX = 0
	}
	if mouse.IsButtonJustPressed(btn.Left) && focused {
		if typingIn != widget {
			scrollX = 0
		}

		typingIn = widget
		widget.textBox.Text = text
	}
	if typingIn == widget {
		text = tryInput(text, widget, margin, root, cam)
		tryRemove(cam, text, root, widget, margin)
		tryMoveCursor(widget, text, cam, margin, root)
		tryFocusNextField(cam, root, widget)

		scrollX = condition.If(txt.Length(text) == 0, 0, scrollX)
		cursorTime += internal.DeltaTime
	}

	var isPlaceholder = false
	if text == "" {
		var placeholder = root.themedField(field.InputFieldPlaceholder, owner, widget)
		placeholder = txt.Remove(defaultValue(placeholder, "Type..."), "\n")
		setupText(margin, root, widget, false) // don't skip when empty!
		widget.textBox.Text = placeholder
		isPlaceholder = true
	}

	draw(margin, cam, root, widget, isPlaceholder)
	cursorTime = condition.If(cursorTime > 1, 0, cursorTime)
}
func draw(margin float32, cam *graphics.Camera, root *root, widget *widget, isPlaceholder bool) {
	maskText = true
	textMargin = margin
	drawVisuals(cam, root, widget, isPlaceholder, func() {
		if indexCursor == indexSelect || typingIn != widget || len(symbolXs) == 0 {
			return
		}
		var ax = symbolXs[indexCursor] - scrollX
		var bx = symbolXs[indexSelect] - scrollX

		if ax > bx {
			ax, bx = bx, ax
		}

		var y, h = widget.textBox.Y + margin/2, widget.textBox.Height - margin
		cam.DrawQuad(ax, y, bx-ax, h, 0, palette.Azure)
	})
	maskText = false

	if typingIn == widget {
		cam.DrawQuadFrame(widget.X, widget.Y, widget.Width, widget.Height, 0, -4, palette.Gray)
	}

	if typingIn == widget && cursorTime < 0.5 {
		var x, y = cursorX(margin, widget), widget.textBox.Y + margin/2
		var w, h = cursorWidth, widget.textBox.Height - margin
		cam.DrawQuad(x, y, w, h, 0, palette.White)
		cam.DrawQuadFrame(x, y, w, h, 0, w/2, palette.Black)
	}
}

func tryMoveCursor(widget *widget, text string, cam *graphics.Camera, margin float32, root *root) {
	var ctrl = keyboard.IsKeyPressed(key.LeftControl) || keyboard.IsKeyPressed(key.RightControl)
	var home = keyboard.IsKeyJustPressed(key.UpArrow) || keyboard.IsKeyJustPressed(key.Home)
	var end = keyboard.IsKeyJustPressed(key.DownArrow) || keyboard.IsKeyJustPressed(key.End)
	var length = txt.Length(text)

	if keyboard.IsKeyJustPressed(key.LeftArrow) || keyboard.IsKeyHeld(key.LeftArrow) {
		var max = number.Biggest(indexCursor-1, 0)
		cursorTime = 0
		indexCursor = condition.If(ctrl, wordIndex(text, true), max)
		trySelect()
	}
	if keyboard.IsKeyJustPressed(key.RightArrow) || keyboard.IsKeyHeld(key.RightArrow) {
		var min = number.Smallest(length, indexCursor+1)
		cursorTime = 0
		indexCursor = condition.If(ctrl, wordIndex(text, false), min)
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
			indexCursor = closestIndexToMouse(cam)

			if mouse.IsButtonJustPressed(btn.Left) {
				calculateXs(widget, cam) // calculate once and update indexes to not drop performance
				indexCursor = closestIndexToMouse(cam)
				indexSelect = indexCursor
			}
		}
	}

	var cx = cursorX(margin, widget)
	var left, right = widget.X + margin, widget.X + widget.Width - margin
	if cx < left && indexCursor >= 0 {
		scrollX -= left - cx
		setupText(margin, root, widget, true)
	}
	if cx > right && indexCursor <= length {
		scrollX += cx - right
		setupText(margin, root, widget, true)
	}
}
func closestIndexToMouse(cam *graphics.Camera) int {
	var mx, _ = cam.MousePosition()
	mx += scrollX

	if len(symbolXs) == 0 {
		return 0
	}

	var closestIndex = 0
	var minDist = number.Unsign(mx - symbolXs[0])

	for i, v := range symbolXs[1:] {
		var dist = number.Unsign(mx - v)
		if dist < minDist {
			minDist = dist
			closestIndex = i + 1
		}
	}

	return closestIndex
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
func tryRemove(cam *graphics.Camera, text string, root *root, widget *widget, margin float32) {
	var left, right = widget.X + margin, widget.X + widget.Width - margin
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
		setText(widget, text)
		indexCursor -= back
		indexSelect = indexCursor
		calculateXs(widget, cam)

		var owner = root.Containers[widget.OwnerId]
		sound.AssetId = defaultValue(root.themedField(field.InputFieldSoundErase, owner, widget), "~erase")
		sound.Volume = root.Volume
		sound.Play()
	}

	if keyboard.IsKeyJustPressed(key.Backspace) || keyboard.IsKeyJustPressed(key.Delete) || simulateRemove {
		if indexSelect < indexCursor {
			remove(indexCursor-indexSelect, 0)
			return
		} else if indexCursor < indexSelect {
			remove(0, indexSelect-indexCursor)
			return
		}
	}

	if keyboard.IsKeyJustPressed(key.Backspace) || keyboard.IsKeyHeld(key.Backspace) {
		remove(condition.If(ctrl, indexCursor-wordIndex(text, true), 1), 0)

		// scrolls left when empty space appears on the right (if possible)
		var textWidth, _ = widget.textBox.TextMeasure(widget.textBox.Text)
		var textRight = (left - scrollX) + textWidth
		if indexCursor > 0 && textRight < right {
			scrollX -= right - textRight
			scrollX = condition.If(textWidth < right-left, 0, scrollX)
			setupText(margin, root, widget, true)
		}
	}
	if keyboard.IsKeyJustPressed(key.Delete) || keyboard.IsKeyHeld(key.Delete) {
		remove(0, condition.If(ctrl, wordIndex(text, false)-indexCursor, 1))
	}
}
func tryInput(text string, widget *widget, margin float32, root *root, cam *graphics.Camera) string {
	var input = keyboard.Input()
	if input == "" {
		return text
	}

	if indexCursor != indexSelect { // text is selected, we should remove it and then type
		simulateRemove = true
		tryRemove(cam, text, root, widget, margin)
		text = widget.textBox.Text
		simulateRemove = false
	}

	if txt.Length(text) == 0 {
		text = input
		setText(widget, text)
		setupText(margin, root, widget, true) // text is not setuped cuz it was empty "" (skipped)
	} else {
		text = txt.Insert(text, input, indexCursor)
	}

	var owner = root.Containers[widget.OwnerId]
	sound.AssetId = defaultValue(root.themedField(field.InputFieldSoundType, owner, widget), "~write")
	sound.Volume = root.Volume
	sound.Play()

	setText(widget, text)
	indexCursor += txt.Length(input)
	indexSelect = indexCursor
	cursorTime = 0
	calculateXs(widget, cam)
	return text
}
func tryFocusNextField(cam *graphics.Camera, root *root, self *widget) {
	if !keyboard.IsKeyJustPressed(key.Tab) || frame == int(time.FrameCount()) {
		return
	}

	var owner = root.Containers[self.OwnerId]
	var allInputFields = []*widget{}
	var myIndex = 0
	for _, wId := range owner.Widgets {
		var w = root.Widgets[wId]

		if w.Class == "inputField" && !w.IsCulled && !w.isHidden(root, owner) {
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
	var text = txt.Remove(root.themedField(field.Text, owner, typingIn), "\n")
	indexCursor = len(text)
	indexSelect = indexCursor
	frame = int(time.FrameCount()) // only once per frame

	var margin = parseNum(root.themedField(field.InputFieldMargin, owner, typingIn), 10)
	setupText(margin, root, typingIn, true)
	if text == "" { // empty text is skipped in setupText so Xs should affect that
		typingIn.textBox.Text = ""
	}
	calculateXs(typingIn, cam)
}

func calculateXs(self *widget, cam *graphics.Camera) {
	var textLength = txt.Length(self.textBox.Text)
	symbolXs = []float32{}

	for i := range textLength {
		var x, _, _, _, _ = self.textBox.TextSymbol(cam, i)
		symbolXs = append(symbolXs, x+scrollX)
	}
	if len(symbolXs) > 0 {
		var w, _ = self.textBox.TextMeasure(self.textBox.Text)
		symbolXs = append(symbolXs, self.textBox.X+w+scrollX)
	}

	if indexSelect > textLength {
		indexSelect = textLength
	}
}
func cursorX(margin float32, widget *widget) float32 {
	var length = len(symbolXs)
	if length > 0 && indexCursor < length {
		return symbolXs[indexCursor] - scrollX
	}
	return widget.X + margin
}
func wordIndex(text string, left bool) int {
	if left && text == "" {
		return 0
	}

	var length = txt.Length(text)
	if !left && indexCursor == length {
		return length
	}

	var symbolIndex = number.Limit(indexCursor, 0, length-1)
	if left {
		symbolIndex = number.Limit(indexCursor-1, 0, length-1)
	}

	var isSpace = text[symbolIndex] == ' '

	if left {
		for i := symbolIndex; i >= 0; i-- {
			if (isSpace && text[i] != ' ') || (!isSpace && text[i] == ' ') {
				return i + 1
			}
		}
		return 0
	}

	for i := indexCursor; i < length; i++ {
		if (isSpace && text[i] != ' ') || (!isSpace && text[i] == ' ') {
			return i
		}
	}
	return length

}
func setText(widget *widget, text string) {
	widget.Fields[field.Text] = text
	widget.textBox.Text = text
}
