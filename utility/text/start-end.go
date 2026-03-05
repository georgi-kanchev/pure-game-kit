package text

import (
	"fmt"
	"strings"
)

func Pad(text string, length int, pad string) string {
	var textLen = Length(text)
	var spaces = length - textLen
	if spaces <= 0 {
		return text
	}
	var left = spaces / 2
	return PadRight(PadLeft(text, textLen+left, pad), length, pad)
}
func PadLeft(text string, length int, pad string) string {
	var textLen = Length(text)
	var padding = length - textLen
	if padding <= 0 || pad == "" {
		return text
	}
	return repeatPad(pad, padding) + text
}
func PadRight(text string, length int, pad string) string {
	var textLen = Length(text)
	var padding = length - textLen
	if padding <= 0 || pad == "" {
		return text
	}
	return text + repeatPad(pad, padding)
}
func PadZeros(number float32, amountOfZeros int) string {
	if amountOfZeros == 0 {
		return New(number)
	}
	if amountOfZeros < 0 {
		var width = -amountOfZeros
		return fmt.Sprintf("%0*d", width, int(number))
	}
	return fmt.Sprintf("%.*f", amountOfZeros, number)
}

func Trim(text string) string {
	return TrimEnd(TrimStart(text))
}
func TrimStart(text string) string {
	return strings.TrimLeft(text, " \r\n")
}
func TrimEnd(text string) string {
	return strings.TrimRight(text, " \r\n")
}

// Surrounds a text with the given start part and end part.
// If end part is empty, it uses the start part for both sides.
func SurroundWith(text, startPart string, endPart ...string) string {
	var end = startPart
	if len(endPart) > 0 {
		end = endPart[0]
	}
	return startPart + text + end
}

// Adds the part to the start of the text only if it doesn't already have it.
func EnsureStart(text, part string) string {
	if !strings.HasPrefix(text, part) {
		return part + text
	}
	return text
}

// Adds the part to the end of the text only if it doesn't already have it.
func EnsureEnd(text, part string) string {
	if !strings.HasSuffix(text, part) {
		return text + part
	}
	return text
}

// Removes the given start part and end part only if both are present.
// If end part is empty, it looks for the start part on both sides.
func Chop(text, startPart string, endPart ...string) string {
	var end = startPart
	if len(endPart) > 0 {
		end = endPart[0]
	}

	if strings.HasPrefix(text, startPart) && strings.HasSuffix(text, end) {
		return text[len(startPart) : len(text)-len(end)]
	}
	return text
}

// Removes the start part from the text if it exists.
func ChopStart(text, part string) string {
	return strings.TrimPrefix(text, part)
}

// Removes the end part from the text if it exists.
func ChopEnd(text, part string) string {
	return strings.TrimSuffix(text, part)
}
