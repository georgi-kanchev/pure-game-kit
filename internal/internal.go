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

func ReplaceStrings(text string, open, close, placeholder rune) (replaced string, originals []string) {
	var result = txt.NewBuilder()
	var current = txt.NewBuilder()
	var inside = false

	for i, char := range text {
		if char == open && !inside {
			inside = true
			continue
		}

		if char == close && inside {
			inside = false
			originals = append(originals, current.ToText())
			result.WriteSymbol(placeholder)
			current.Clear()
			continue
		}

		if inside {
			current.WriteSymbol(char)
			if i == len(text)-1 {
				result.WriteSymbol(open)
				result.WriteText(current.ToText())
			}
		} else {
			result.WriteSymbol(char)
		}
	}

	return result.ToText(), originals
}
