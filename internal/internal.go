// Used by the engine packages only:
//   - As a third party messenger when communicating with each other, avoiding dependencies loops.
//   - To re-use private engine (non-API) code.
//   - To store runtime resources such as assets, callbacks, input cache, time cache and so on.
//   - To pump updates every frame coming from the window package onto other packages that require them.
package internal

import (
	"regexp"
	"strings"
)

var tagRegexp = regexp.MustCompile(`{.*?}`)
var tagBuffer strings.Builder

const Placeholder = '╌'

func ReplaceStrings(text string, open, close, placeholder rune) (replaced string, originals []string) {
	var result strings.Builder
	var current strings.Builder
	var inside = false

	for i, char := range text {
		if char == open && !inside {
			inside = true
			continue
		}

		if char == close && inside {
			inside = false
			originals = append(originals, current.String())
			result.WriteRune(placeholder)
			current.Reset()
			continue
		}

		if inside {
			current.WriteRune(char)
			if i == len(text)-1 {
				result.WriteRune(open)
				result.WriteString(current.String())
			}
		} else {
			result.WriteRune(char)
		}
	}

	return result.String(), originals
}
func RemoveTags(text string) string {
	if !strings.ContainsRune(text, '<') {
		return text // FAST PATH: If there are no brackets, return the original string - ZERO allocations for tagless text
	}

	tagBuffer.Reset()
	tagBuffer.Grow(len(text)) // pre-size the buffer to the length of the text to avoid intermediate grows

	var inTag = false
	for i := 0; i < len(text); i++ {
		var char = text[i]
		if char == '<' {
			inTag = true
			continue
		}
		if char == '>' {
			inTag = false
			continue
		}
		if !inTag {
			tagBuffer.WriteByte(char)
		}
	}

	return tagBuffer.String()
}
