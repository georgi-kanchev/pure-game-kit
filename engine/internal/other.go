package internal

import (
	"strings"
)

const Placeholder = '╌'

func ReplaceQuotedStrings(text string, quote, placeholder rune) (replaced string, originals []string) {
	var result strings.Builder
	var inQuotes bool
	var current strings.Builder

	for i, char := range text {
		if char == quote {
			if inQuotes { // closing quote found
				inQuotes = false
				originals = append(originals, current.String())
				result.WriteRune(placeholder)
				current.Reset()
			} else { // opening quote found
				inQuotes = true
			}
			continue
		}

		if inQuotes {
			current.WriteRune(char)

			if i == len(text)-1 { // no closing quote found results in no replacement
				result.WriteString(string(quote) + current.String())
			}
		} else {
			result.WriteRune(char)
		}
	}

	return result.String(), originals
}
