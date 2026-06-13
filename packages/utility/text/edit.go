package text

import (
	"fmt"
	"pure-game-kit/packages/utility/number"
	"strings"
)

func Insert(text, part string, atIndex int) string {
	var lastIndex = Length(text)
	var before = Part(text, 0, atIndex)
	var after = Part(text, atIndex, lastIndex)
	return before + part + after
}
func Part(text string, fromIndex, toIndex int) string {
	var runes = []rune(text)
	var length = len(runes)
	var start = number.Limit(fromIndex, 0, length)
	var end = number.Limit(toIndex, 0, length)

	if start > end {
		return ""
	}
	return string(runes[start:end])
}
func Replace(text, part, with string) string {
	return strings.ReplaceAll(text, part, with)
}
func Remove(text string, parts ...string) string {
	for _, part := range parts {
		text = Replace(text, part, "")
	}
	return text
}

// Progress 0..1 for start-to-end and 0..-1 for end-to-start.
func Reveal(text string, progress float32) string {
	var textLen = float32(Length(text))
	if progress >= 0 {
		progress = number.Limit(progress, 0, 1)
		var cutoff = int(number.Round(progress * textLen))
		return string([]rune(text)[cutoff:])
	}
	progress = number.Limit(progress, -1, 0)
	var cutoff = int(number.Round(progress * textLen))
	return string([]rune(text)[cutoff:])
}

// Positive length trims from the end, negative length trims from the start. Default indicator if skipped: '…'
func Limit(text string, length int, indicator ...string) string {
	if length == 0 {
		return ""
	}

	var ind = "…"
	if len(indicator) > 0 {
		ind = indicator[0]
	}
	var textRunes = []rune(text)
	var indicatorLen = len([]rune(ind))
	var textLen = len(textRunes)
	var absMax = number.Unsign(length)
	var trimLen = absMax - indicatorLen

	if length > 0 && textLen > int(length) {
		if trimLen <= 0 {
			return ind
		}
		return string(textRunes[:trimLen]) + ind
	} else if length < 0 && textLen > absMax {
		if trimLen <= 0 {
			return ind
		}
		return ind + string(textRunes[textLen-trimLen:])
	}
	return text
}

// Returns the text before the part. If the part is not found, returns the original text.
func Before(text, part string) string {
	var before, _, ok = strings.Cut(text, part)
	if !ok {
		return text
	}
	return before
}

// Returns the text after the part. If the part is not found, returns the original text.
func After(text, part string) string {
	var _, after, ok = strings.Cut(text, part)
	if !ok {
		return text
	}
	return after
}

// Returns the text between the first part and the second part. If one of the parts is not found,
// returns the original text.
func Between(text, firstPart, secondPart string) string {
	var _, after, ok = strings.Cut(text, firstPart)
	if !ok {
		return text
	}

	var remaining = after // start searching for the endAnchor AFTER the startAnchor
	var ePos = strings.Index(remaining, secondPart)
	if ePos == -1 {
		return text
	}
	return remaining[:ePos]
}

// Returns a text consisting of the provided amount of copies of the text. If count is 0 or negative, returns "".
func Repeat(text string, count int) string {
	if count <= 0 {
		return ""
	}
	return strings.Repeat(text, count)
}

// Breaks the text at the provided line length, regardless of word boundaries.
func Wrap(text string, lineLength int) string {
	if lineLength <= 0 || len(text) <= lineLength {
		return text
	}

	builder.Reset()
	var runes = []rune(text)
	for i, r := range runes {
		builder.WriteRune(r)
		var reachedEnd = (i+1)%lineLength == 0
		var notLastChar = i+1 < len(runes)
		if reachedEnd && notLastChar {
			builder.WriteByte('\n')
		}
	}
	return builder.String()
}

// Breaks the text into lines with a maximum length but only breaks at whitespace to keep words intact.
func WrapWords(text string, lineLength int) string {
	if lineLength <= 0 || len(text) <= lineLength {
		return text
	}

	var words = strings.Fields(text)
	if len(words) == 0 {
		return ""
	}

	builder.Reset()
	var currentLineLength = 0
	for i, word := range words {
		var wordLen = len([]rune(word))
		if currentLineLength+wordLen > lineLength && currentLineLength > 0 { // longer than line length
			builder.WriteByte('\n')
			currentLineLength = 0
		}
		if currentLineLength > 0 {
			builder.WriteByte(' ') // add a space if it's not the start of a new line
			currentLineLength++
		}

		builder.WriteString(word)
		currentLineLength += wordLen

		if currentLineLength >= lineLength && i < len(words)-1 { // still longer than line length
			builder.WriteByte('\n') // the next word MUST start on a new line
			currentLineLength = 0
		}
	}
	return builder.String()
}

// Converts the text to lowercase
func ToLowerCase(text string) string {
	return strings.ToLower(text)
}

// Converts the text to UPPERCASE
func ToUpperCase(text string) string {
	return strings.ToUpper(text)
}

func FormatByteSize(byteSize int) string {
	const unit = 1024
	if byteSize < unit {
		return fmt.Sprintf("%d B", byteSize)
	}
	var div, exp = int(unit), 0
	for n := byteSize / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.3f %cB", float32(byteSize)/float32(div), "KMGTPE"[exp])
}
