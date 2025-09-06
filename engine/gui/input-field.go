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

func setupText(margin float32, root *root, widget *widget, owner *container) {
	setupVisualsText(root, widget, owner)
	reusableTextBox.AlignmentX, reusableTextBox.AlignmentY = 0, 0
	reusableTextBox.Width = 9999
	reusableTextBox.Height = widget.Height - margin
	reusableTextBox.LineHeight = widget.Height - margin
	reusableTextBox.X = reusableTextBox.X + margin - scrollX
	reusableTextBox.Y += margin / 2
	reusableTextBox.Text = strings.ReplaceAll(reusableTextBox.Text, "\n", "")
	reusableTextBox.EmbeddedAssetsTag = 0
	reusableTextBox.EmbeddedColorsTag = 0
	reusableTextBox.EmbeddedThicknessesTag = 0
}
func inputField(cam *graphics.Camera, root *root, widget *widget, owner *container) {
	var margin = parseNum(themedProp(property.InputFieldMargin, root, owner, widget), 30)
	setupText(margin, root, widget, owner)
	setupVisualsTextured(root, widget, owner)

	if widget.isFocused(root, cam) {
		mouse.SetCursor(mouse.CursorInput)
	}

	var focused = widget.isFocused(root, cam)
	var anyInput = mouse.IsAnyButtonPressedOnce() || mouse.Scroll() != 0

	if anyInput && !focused {
		typingIn = nil
		scrollX = 0
	}
	if mouse.IsButtonPressedOnce(mouse.ButtonLeft) && focused {
		typingIn = widget
		cursorTime = 0
		var text = strings.ReplaceAll(themedProp(property.Text, root, owner, widget), "\n", "")
		if txt.Count(text) == 0 {
			symbolXs = []float32{}
			indexCursor = 0
		} else {
			calculateXs(cam)
			indexCursor = closestIndex(cam)
		}
	}
	if typingIn == widget {
		cam.DrawFrame(widget.X, widget.Y, widget.Width, widget.Height, 0, -5, color.Gray)

		var input = keyboard.Input()
		var text = strings.ReplaceAll(themedProp(property.Text, root, owner, widget), "\n", "")
		var left, right = widget.X + margin, widget.X + widget.Width - margin
		if input != "" {
			if txt.Count(text) == 0 {
				text = input
				widget.Properties[property.Text] = text
				setupText(margin, root, widget, owner) // text is not setuped cuz it was empty "" (skipped)
			} else {
				text = text[:indexCursor] + input + text[indexCursor:]
			}
			widget.Properties[property.Text] = text
			reusableTextBox.Text = text
			indexCursor += txt.Count(input)
			cursorTime = 0
			calculateXs(cam)
		}
		if keyboard.IsKeyPressedOnce(key.Backspace) || keyboard.IsKeyHeld(key.Backspace) {
			remove(1, 0, cam, root, widget, owner)

			var textWidth, _ = reusableTextBox.TextMeasure(reusableTextBox.Text)
			var textRight = (left - scrollX) + textWidth
			if indexCursor > 0 && textRight < right {
				scrollX -= right - textRight
				scrollX = condition.If(textWidth < right-left, 0, scrollX)
				setupText(margin, root, widget, owner)
			}
		}
		if keyboard.IsKeyPressedOnce(key.Delete) || keyboard.IsKeyHeld(key.Delete) {
			remove(0, 1, cam, root, widget, owner)
		}

		if keyboard.IsKeyPressedOnce(key.Left) || keyboard.IsKeyHeld(key.Left) {
			indexCursor = number.BiggestInt(indexCursor-1, 0)
			cursorTime = 0

		}
		if keyboard.IsKeyPressedOnce(key.Right) || keyboard.IsKeyHeld(key.Right) {
			indexCursor = number.SmallestInt(txt.Count(text), indexCursor+1)
			cursorTime = 0
		}

		var cx = cursorX(margin, widget)
		if cx < left && indexCursor >= 0 {
			scrollX -= left - cx
			setupText(margin, root, widget, owner)
		}
		if cx > right && indexCursor <= txt.Count(text) {
			scrollX += cx - right
			setupText(margin, root, widget, owner)
		}

		scrollX = condition.If(txt.Count(text) == 0, 0, scrollX)
		cursorTime += seconds.RealFrameDelta()
	}

	maskText = true
	textMargin = margin
	drawVisuals(cam, root, widget, owner)
	maskText = false

	if typingIn == widget && cursorTime < 0.5 {
		var x = cursorX(margin, widget)
		cam.DrawLine(x, reusableTextBox.Y, x, reusableTextBox.Y+reusableTextBox.LineHeight, 5, color.Black)
	}
	cursorTime = condition.If(cursorTime > 1, 0, cursorTime)
}

func closestIndex(cam *graphics.Camera) int {
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
func calculateXs(cam *graphics.Camera) {
	var textLength = txt.Count(reusableTextBox.Text)
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
func remove(back, front int, cam *graphics.Camera, root *root, widget *widget, owner *container) {
	cursorTime = 0

	if back > 0 && indexCursor == 0 {
		return
	}
	var text = strings.ReplaceAll(themedProp(property.Text, root, owner, widget), "\n", "")
	if front > 0 && indexCursor == txt.Count(text) {
		return
	}
	text = text[:indexCursor-back] + text[indexCursor+front:]
	widget.Properties[property.Text] = text
	reusableTextBox.Text = text
	indexCursor -= back
	calculateXs(cam)
}
