package gui

import (
	"pure-game-kit/graphics"
	"pure-game-kit/gui/field"
	"pure-game-kit/utility/number"
	txt "pure-game-kit/utility/text"
	"unicode"
)

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

func wordIndex(text string, left bool, fromIndex int) int {
	var runes = []rune(text)
	var length = len(runes)
	var isLetterOrDigit = func(r rune) bool {
		return unicode.IsLetter(r) || unicode.IsDigit(r)
	}

	if length == 0 {
		return 0
	}

	var cursor = number.Limit(fromIndex, 0, length)

	if left {
		if cursor == 0 {
			return 0
		}
		var startIdx = cursor - 1
		var isStartWord = isLetterOrDigit(runes[startIdx])

		for i := startIdx; i >= 0; i-- {
			if isLetterOrDigit(runes[i]) != isStartWord {
				return i + 1
			}
		}
		return 0
	}

	if cursor >= length {
		return length
	}

	var isStartWord = isLetterOrDigit(runes[cursor])
	for i := cursor; i < length; i++ {
		if isLetterOrDigit(runes[i]) != isStartWord {
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
