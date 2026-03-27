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
	"pure-game-kit/utility/is"
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

func setupText(margin float32, w *widget, skipEmpty bool) {
	setupVisualsText(w, skipEmpty)
	w.textBox.AlignmentX, w.textBox.AlignmentY = 0, 0
	w.textBox.Width = 99999
	w.textBox.Height = w.Height - margin
	w.textBox.LineHeight = w.Height - margin
	var scroll = condition.If(typingIn == w, scrollX, 0)
	w.textBox.X = w.textBox.X + margin - scroll
	w.textBox.Y += margin / 2
	w.textBox.Text = txt.Remove(w.textBox.Text, "\n")
}
func inputField(w *widget) {
	maskText = true
	if w.textBox == nil {
		w.textBox = &graphics.TextBox{}
	}

	var owner = w.root.Containers[w.OwnerId]
	var margin = parseNum(w.root.themedField(field.InputFieldMargin, owner, w), 10)

	if keyboard.IsComboJustPressed(key.LeftControl, key.A) ||
		keyboard.IsComboJustPressed(key.RightControl, key.A) {
		indexSelect = 0
		indexCursor = len(symbolXs) - 1
	}

	setupText(margin, w, true)
	setupVisualsTextured(w)

	if w.isFocused() {
		mouse.SetCursor(cursor.Input)
	}

	var anyInput = mouse.IsAnyButtonJustPressed() || mouse.Scroll() != 0
	var focused = w.isFocused()
	var meTyping = typingIn == w // each input field should disable its own typing
	var text = txt.Remove(w.root.themedField(field.Text, owner, w), "\n")
	text = internal.RemoveTags(text)

	if meTyping && ((anyInput && !focused) || !window.IsHovered() || keyboard.IsKeyJustPressed(key.Escape)) {
		typingIn = nil
		scrollX = 0
	}
	if mouse.IsButtonJustPressed(btn.Left) && focused {
		if typingIn != w {
			scrollX = 0
		}

		typingIn = w
		w.textBox.Text = text
	}
	if typingIn == w {
		text = tryInput(text, w, margin)
		text = tryRemove(text, w, margin)
		tryMoveCursor(w, text, margin)
		tryFocusNextField(w)

		scrollX = condition.If(txt.Length(text) == 0, 0, scrollX)
		cursorTime += internal.DeltaTime
	}

	var isPlaceholder = false
	if text == "" {
		var placeholder = w.root.themedField(field.InputFieldPlaceholder, owner, w)
		placeholder = txt.Remove(defaultValue(placeholder, "Type..."), "\n")
		setupText(margin, w, false) // don't skip when empty!
		w.textBox.Text = placeholder
		isPlaceholder = true
	}

	draw(margin, w, isPlaceholder)
	cursorTime = condition.If(cursorTime > 1, 0, cursorTime)
	maskText = false
}
func draw(margin float32, w *widget, isPlaceholder bool) {
	if typingIn == w {
		if w.highlight == nil {
			w.highlight = graphics.NewBox("", 0, 0)
		}
		w.highlight.X, w.highlight.Y = w.X-0.5, w.Y-0.5
		w.highlight.Width, w.highlight.Height = w.Width+0.5, w.Height+0.5
		w.highlight.Tint = palette.Azure
		w.highlight.PivotX, w.highlight.PivotY = 0, 0
		w.highlight.Mask = w.root.Containers[w.OwnerId].mask
		w.root.boxes = append(w.root.boxes, w.highlight)
	}

	textMargin = margin
	drawVisuals(w, isPlaceholder, func() {
		if indexCursor == indexSelect || typingIn != w || len(symbolXs) == 0 {
			return
		}
		var ax = symbolXs[indexCursor] - scrollX
		var bx = symbolXs[indexSelect] - scrollX

		if ax > bx {
			ax, bx = bx, ax
		}

		if w.handle == nil {
			w.handle = graphics.NewSprite("", 0, 0)
		}

		w.handle.X, w.handle.Y = ax, w.textBox.Y+margin/2
		w.handle.Width, w.handle.Height = bx-ax, w.textBox.Height-margin
		w.handle.Tint, w.handle.Mask = palette.Azure, w.textBox.Mask
		w.handle.PivotX, w.handle.PivotY = 0, 0
		w.root.sprites = append(w.root.sprites, w.handle)
	})

	if typingIn == w && cursorTime < 0.5 {
		var x, y = cursorX(margin, w), w.textBox.Y + margin/2
		var cw, ch = cursorWidth, w.textBox.Height - margin

		if w.cursor1 == nil || w.cursor2 == nil {
			w.cursor1 = graphics.NewSprite("", 0, 0)
			w.cursor2 = graphics.NewSprite("", 0, 0)
		}

		w.cursor1.X, w.cursor1.Y = x-cw/2, y-cw/2
		w.cursor1.Width, w.cursor1.Height = cw+cw, ch+cw
		w.cursor1.Tint, w.cursor1.Mask = palette.Azure, w.textBox.Mask
		w.cursor1.PivotX, w.cursor1.PivotY = 0, 0

		w.cursor2.X, w.cursor2.Y = x, y
		w.cursor2.Width, w.cursor2.Height = cw, ch
		w.cursor2.Tint, w.cursor2.Mask = palette.Black, w.textBox.Mask
		w.cursor2.PivotX, w.cursor2.PivotY = 0, 0
		w.root.spritesAbove = append(w.root.spritesAbove, w.cursor1, w.cursor2)
	}
}

func tryMoveCursor(w *widget, text string, margin float32) {
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
func tryRemove(text string, w *widget, margin float32) string {
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
		remove(condition.If(ctrl, indexCursor-wordIndex(text, true), 1), 0)

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
		remove(0, condition.If(ctrl, wordIndex(text, false)-indexCursor, 1))
	}
	return text
}
func tryInput(text string, w *widget, margin float32) string {
	var input = keyboard.Input()
	if is.AnyOf(input, "", "{", "}") {
		return text
	}

	if indexCursor != indexSelect { // text is selected, we should remove it and then type
		simulateRemove = true
		text = tryRemove(text, w, margin)
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
func tryFocusNextField(self *widget) {
	if !keyboard.IsKeyJustPressed(key.Tab) || frame == int(time.FrameCount()) {
		return
	}

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

func calculateXs(self *widget) {
	var textLength = txt.Length(self.textBox.Text)
	symbolXs = []float32{}

	for i := range textLength {
		var x, _, _, _, _ = self.textBox.TextSymbol(i)
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
func cursorX(margin float32, w *widget) float32 {
	var length = len(symbolXs)
	if length > 0 && indexCursor < length {
		return symbolXs[indexCursor] - scrollX
	}
	return w.X + margin
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
	text = txt.Remove(text, "{", "}")
	widget.Fields[field.Text] = text
	widget.textBox.Text = text
}
