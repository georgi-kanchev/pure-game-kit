package internal

import txt "pure-kit/engine/utility/text"

const Placeholder = '╌'

func ReplaceQuotedStrings(text string, quote, placeholder rune) (replaced string, originals []string) {
	var result = txt.NewBuilder()
	var inQuotes bool
	var current = txt.NewBuilder()

	for i, char := range text {
		if char == quote {
			if inQuotes { // closing quote found
				inQuotes = false
				originals = append(originals, current.ToText())
				result.WriteSymbol(placeholder)
				current.Clear()
			} else { // opening quote found
				inQuotes = true
			}
			continue
		}

		if inQuotes {
			current.WriteSymbol(char)

			if i == len(text)-1 { // no closing quote found results in no replacement
				result.WriteText(string(quote) + current.ToText())
			}
		} else {
			result.WriteSymbol(char)
		}
	}

	return result.ToText(), originals
}
