package gui

import (
	"math"
	"pure-kit/engine/execution/condition"
	"pure-kit/engine/graphics"
	"pure-kit/engine/gui/property"
	"pure-kit/engine/input/mouse"
	"pure-kit/engine/utility/color"
	"pure-kit/engine/utility/seconds"
	"pure-kit/engine/utility/symbols"
	"strings"
)

func InputField(id string, properties ...string) string {
	return newWidget("inputField", id, properties...)
}

// #region private

var typingIn *widget
var indexCursor, indexSelect int
var cursorTime float32
var symbolXs []float32 = []float32{}

func inputField(cam *graphics.Camera, root *root, widget *widget, owner *container) {
	setupVisualsText(root, widget, owner)
	setupVisualsTextured(root, widget, owner)

	var margin = parseNum(themedProp(property.InputFieldMargin, root, owner, widget), 30)
	reusableTextBox.AlignmentX, reusableTextBox.AlignmentY = 0, 0.5
	reusableTextBox.LineHeight = widget.Height - margin
	reusableTextBox.X += margin
	reusableTextBox.Width -= margin * 2
	reusableTextBox.Text = strings.ReplaceAll(reusableTextBox.Text, "\n", "")

	if widget.IsFocused(root, cam) {
		mouse.SetCursor(mouse.CursorInput)
	}

	drawVisuals(cam, root, widget, owner)

	var focused = widget.IsFocused(root, cam)
	var anyInput = mouse.IsAnyButtonPressedOnce() || mouse.Scroll() != 0

	if anyInput && !focused {
		typingIn = nil
	}
	if mouse.IsButtonPressedOnce(mouse.ButtonLeft) && focused {
		var textLength = symbols.Count(reusableTextBox.Text)

		typingIn = widget
		cursorTime = 0
		symbolXs = []float32{}

		for i := range textLength {
			var x, _, _, _, _ = reusableTextBox.TextSymbol(cam, i)
			symbolXs = append(symbolXs, x)
		}
		if len(symbolXs) > 0 {
			var w, _ = reusableTextBox.TextMeasure(reusableTextBox.Text)
			symbolXs = append(symbolXs, reusableTextBox.X+w)
		}
		indexCursor = closestIndex(cam)

	}
	if typingIn == widget {
		cam.DrawFrame(widget.X, widget.Y, widget.Width, widget.Height, 0, -3, color.Gray)

		cursorTime += seconds.RealFrameDelta()

		if cursorTime < 0.5 {
			var x, y = symbolXs[indexCursor], reusableTextBox.Y + margin/2
			cam.DrawLine(x, y, x, y+reusableTextBox.LineHeight, 5, color.Black)
		}
	}

	cursorTime = condition.If(cursorTime > 1, 0, cursorTime)
}

func closestIndex(cam *graphics.Camera) int {
	var mx, _ = cam.MousePosition()

	if len(symbolXs) == 0 {
		return 0
	}

	var closestIndex = 0
	var minDist = float32(math.Abs(float64(mx - symbolXs[0])))

	for i, v := range symbolXs[1:] {
		var dist = float32(math.Abs(float64(mx - v)))
		if dist < minDist {
			minDist = dist
			closestIndex = i + 1
		}
	}

	return closestIndex
}

// #endregion
