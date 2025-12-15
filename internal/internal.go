/*
Used by the engine packages only:
  - As a third party messenger when communicating with each other, avoiding dependencies loops.
  - To re-use private engine (non-API) code.
  - To store runtime resources such as assets, callbacks, input cache, time cache and so on.
  - To pump updates every frame coming from the window package onto other packages that require them.
*/
package internal

import txt "pure-game-kit/utility/text"

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
