package gui

import (
	"pure-kit/engine/execution/condition"
	"pure-kit/engine/graphics"
	"pure-kit/engine/gui/property"
	"pure-kit/engine/input/keyboard"
	"pure-kit/engine/input/keyboard/key"
	"pure-kit/engine/input/mouse"
	"pure-kit/engine/utility/color"
	"pure-kit/engine/utility/number"
	"pure-kit/engine/utility/seconds"
	txt "pure-kit/engine/utility/text"
	"pure-kit/engine/window"
	"strings"
)

func InputField(id string, properties ...string) string {
	return newWidget("inputField", id, properties...)
}

//=================================================================
// private

var typingIn *widget
var indexCursor int
var cursorTime, scrollX, textMargin float32
var symbolXs []float32 = []float32{}
var maskText = false
var frame = 0

func setupText(margin float32, root *root, widget *widget) {
	setupVisualsText(root, widget)
	reusableTextBox.AlignmentX, reusableTextBox.AlignmentY = 0, 0
	reusableTextBox.Width = 9999
	reusableTextBox.Height = widget.Height - margin
	reusableTextBox.LineHeight = widget.Height - margin
	var scroll = condition.If(typingIn == widget, scrollX, 0)
	reusableTextBox.X = reusableTextBox.X + margin - scroll
	reusableTextBox.Y += margin / 2
	reusableTextBox.Text = strings.ReplaceAll(reusableTextBox.Text, "\n", "")
	reusableTextBox.EmbeddedAssetsTag = 0
	reusableTextBox.EmbeddedColorsTag = 0
	reusableTextBox.EmbeddedThicknessesTag = 0
}
func inputField(cam *graphics.Camera, root *root, widget *widget) {
	var owner = root.Containers[widget.OwnerId]
	var margin = parseNum(themedProp(property.InputFieldMargin, root, owner, widget), 30)

	setupText(margin, root, widget)
	setupVisualsTextured(root, widget)

	if widget.isFocused(root, cam) {
		mouse.SetCursor(mouse.CursorInput)
	}

	var focused = widget.isFocused(root, cam)
	var anyInput = mouse.IsAnyButtonPressedOnce() || mouse.Scroll() != 0
	if (anyInput && !focused) || !window.IsHovered() {
		typingIn = nil
		scrollX = 0
	}
	if mouse.IsButtonPressedOnce(mouse.ButtonLeft) && focused {
		if typingIn != widget {
			scrollX = 0
		}

		typingIn = widget
	}
	if typingIn == widget {
		var text = strings.ReplaceAll(themedProp(property.Text, root, owner, widget), "\n", "")
		text = tryInput(text, widget, margin, root, cam)
		tryRemove(cam, text, root, widget, margin)
		tryMoveCursor(widget, text, cam, margin, root)
		tryFocusNextField(cam, root, widget)

		scrollX = condition.If(txt.Length(text) == 0, 0, scrollX)
		cursorTime += seconds.RealFrameDelta()
	}

	maskText = true
	textMargin = margin
	drawVisuals(cam, root, widget)
	maskText = false

	if typingIn == widget {
		cam.DrawFrame(widget.X, widget.Y, widget.Width, widget.Height, 0, -5, color.Gray)
	}

	if typingIn == widget && cursorTime < 0.5 {
		var x = cursorX(margin, widget)
		cam.DrawLine(x, reusableTextBox.Y, x, reusableTextBox.Y+reusableTextBox.LineHeight, 5, color.Black)
	}
	cursorTime = condition.If(cursorTime > 1, 0, cursorTime)
}

func tryMoveCursor(widget *widget, text string, cam *graphics.Camera, margin float32, root *root) {
	if keyboard.IsKeyPressedOnce(key.LeftArrow) || keyboard.IsKeyHeld(key.LeftArrow) {
		cursorTime = 0
		if ctrl() {
			indexCursor = wordIndex(text, true)
		} else {
			indexCursor = number.BiggestInt(indexCursor-1, 0)
		}
	}
	if keyboard.IsKeyPressedOnce(key.RightArrow) || keyboard.IsKeyHeld(key.RightArrow) {
		cursorTime = 0
		if ctrl() {
			indexCursor = wordIndex(text, false)
		} else {
			indexCursor = number.SmallestInt(txt.Length(text), indexCursor+1)
		}
	}
	if mouse.IsButtonPressed(mouse.ButtonLeft) {
		cursorTime = 0
		var closestIndex = func(cam *graphics.Camera) int {
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

		if txt.Length(text) == 0 {
			symbolXs = []float32{}
			indexCursor = 0
		} else {
			calculateXs(cam)
			indexCursor = closestIndex(cam)
		}
	}
	if keyboard.IsKeyPressedOnce(key.DownArrow) || keyboard.IsKeyPressedOnce(key.Home) {
		cursorTime = 0
		indexCursor = 0
	}
	if keyboard.IsKeyPressedOnce(key.UpArrow) || keyboard.IsKeyPressedOnce(key.End) {
		cursorTime = 0
		indexCursor = txt.Length(text)
	}

	var cx = cursorX(margin, widget)
	var left, right = widget.X + margin, widget.X + widget.Width - margin
	if cx < left && indexCursor >= 0 {
		scrollX -= left - cx
		setupText(margin, root, widget)
	}
	if cx > right && indexCursor <= txt.Length(text) {
		scrollX += cx - right
		setupText(margin, root, widget)
	}
}
func tryRemove(cam *graphics.Camera, text string, root *root, widget *widget, margin float32) {
	var left, right = widget.X + margin, widget.X + widget.Width - margin
	var remove = func(back, front int) {
		cursorTime = 0

		if back > 0 && indexCursor == 0 {
			return
		}
		if front > 0 && indexCursor == txt.Length(text) {
			return
		}
		text = text[:indexCursor-back] + text[indexCursor+front:]
		widget.Properties[property.Text] = text
		reusableTextBox.Text = text
		indexCursor -= back
		calculateXs(cam)
	}

	if keyboard.IsKeyPressedOnce(key.Backspace) || keyboard.IsKeyHeld(key.Backspace) {
		if ctrl() {
			var newIndex = wordIndex(text, true)
			setText(widget, text[:newIndex]+text[indexCursor:])
			indexCursor = newIndex
			cursorTime = 0
		} else {
			remove(1, 0)

			var textWidth, _ = reusableTextBox.TextMeasure(reusableTextBox.Text)
			var textRight = (left - scrollX) + textWidth
			if indexCursor > 0 && textRight < right {
				scrollX -= right - textRight
				scrollX = condition.If(textWidth < right-left, 0, scrollX)
				setupText(margin, root, widget)
			}
		}
	}
	if keyboard.IsKeyPressedOnce(key.Delete) || keyboard.IsKeyHeld(key.Delete) {
		if ctrl() {
			var newIndex = wordIndex(text, false)
			setText(widget, text[:indexCursor]+text[newIndex:])
			cursorTime = 0
		} else {
			remove(0, 1)
		}
	}
}
func tryInput(text string, widget *widget, margin float32, root *root, cam *graphics.Camera) string {
	var input = keyboard.Input()
	if input == "" {
		return text
	}

	if txt.Length(text) == 0 {
		text = input
		setText(widget, text)
		setupText(margin, root, widget) // text is not setuped cuz it was empty "" (skipped)
	} else {
		text = text[:indexCursor] + input + text[indexCursor:]
	}
	setText(widget, text)
	indexCursor += txt.Length(input)
	cursorTime = 0
	calculateXs(cam)
	return text
}
func tryFocusNextField(cam *graphics.Camera, root *root, self *widget) {
	if !keyboard.IsKeyPressedOnce(key.Tab) || frame == int(seconds.FrameCount()) {
		return
	}

	var owner = root.Containers[self.OwnerId]
	var allInputFields = []*widget{}
	var myIndex = 0
	for _, wId := range owner.Widgets {
		var w = root.Widgets[wId]

		if w.Class == "inputField" {
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
	var text = strings.ReplaceAll(themedProp(property.Text, root, owner, typingIn), "\n", "")
	indexCursor = len(text)
	frame = int(seconds.FrameCount()) // only once per frame

	var margin = parseNum(themedProp(property.InputFieldMargin, root, owner, typingIn), 30)
	setupText(margin, root, typingIn)
	if text == "" { // empty text is skipped in setupText so Xs should affect that
		reusableTextBox.Text = ""
	}
	calculateXs(cam)
}

func calculateXs(cam *graphics.Camera) {
	var textLength = txt.Length(reusableTextBox.Text)
	symbolXs = []float32{}

	for i := range textLength {
		var x, _, _, _, _ = reusableTextBox.TextSymbol(cam, i)
		symbolXs = append(symbolXs, x+scrollX)
	}
	if len(symbolXs) > 0 {
		var w, _ = reusableTextBox.TextMeasure(reusableTextBox.Text)
		symbolXs = append(symbolXs, reusableTextBox.X+w+scrollX)
	}
}
func cursorX(margin float32, widget *widget) float32 {
	if len(symbolXs) > 0 {
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

	var symbolIndex = number.LimitInt(indexCursor, 0, length-1)
	if left {
		symbolIndex = number.LimitInt(indexCursor-1, 0, length-1)
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
	widget.Properties[property.Text] = text
	reusableTextBox.Text = text
}

func ctrl() bool {
	var keys = keyboard.KeysPressed()
	return len(keys) == 2 && (keys[0] == key.LeftControl || keys[0] == key.RightControl)
}
