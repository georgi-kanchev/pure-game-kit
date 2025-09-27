package naming

import (
	"regexp"
	"unicode"

	"pure-kit/engine/utility/collection"
	"pure-kit/engine/utility/flag"
	"pure-kit/engine/utility/naming/casing"
	"pure-kit/engine/utility/number"
	"pure-kit/engine/utility/random"
	txt "pure-kit/engine/utility/text"
)

func Apply(text string, naming int, divider string) string {
	if divider == "" {
		divider = ""
	}

	if naming == casing.Random {
		var result = txt.NewBuilder()
		var randBool = random.Range(0, 1, number.NaN()) < 0.5
		for _, r := range text {
			if randBool {
				result.WriteSymbol(unicode.ToLower(r))
			} else {
				result.WriteSymbol(unicode.ToUpper(r))
			}
		}
		return result.ToText()
	}

	detectedNaming, detectedDivider := Detect(text)
	var words = []string{text}
	if detectedDivider != "" {
		words = txt.Split(text, detectedDivider)
	}

	if len(words) == 1 &&
		divider != "" &&
		(flag.IsOn(detectedNaming, casing.Camel) || flag.IsOn(detectedNaming, casing.Pascal)) {
		words = txt.Split(addDivCamelPascal(words[0], divider), divider)
	}

	for i := range words {
		word := words[i]

		if flag.IsOn(naming, casing.Lower) {
			word = txt.LowerCase(word)
		}

		if flag.IsOn(naming, casing.Upper) {
			word = txt.UpperCase(word)
		}

		if flag.IsOn(naming, casing.Camel) {
			if i == 0 {
				word = txt.LowerCase(word)
			} else {
				word = capitalize(word)
			}
		}

		if flag.IsOn(naming, casing.Pascal) {
			word = capitalize(word)
		}

		if flag.IsOn(naming, casing.Sentence) {
			if i == 0 {
				word = capitalize(word)
			} else {
				word = txt.LowerCase(word)
			}
		}

		if flag.IsOn(naming, casing.Pingpong) {
			var builder = txt.NewBuilder()
			var isUpper = true
			for _, c := range word {
				if isUpper {
					builder.WriteSymbol(unicode.ToUpper(c))
				} else {
					builder.WriteSymbol(unicode.ToLower(c))
				}
				isUpper = !isUpper
			}
			word = builder.ToText()
		}

		if flag.IsOn(naming, casing.Pongping) {
			var builder = txt.NewBuilder()
			var isLower = true
			for _, c := range word {
				if isLower {
					builder.WriteSymbol(unicode.ToLower(c))
				} else {
					builder.WriteSymbol(unicode.ToUpper(c))
				}
				isLower = !isLower
			}
			word = builder.ToText()
		}

		words[i] = word
	}

	return collection.ToText(words, divider)
}
func Detect(text string) (naming int, separator string) {
	if txt.Trim(text) == "" {
		return casing.Random, ""
	}

	var detectedNaming = casing.Random
	var divider = ""
	var words = []string{text}
	var re = regexp.MustCompile(`[^a-zA-Z0-9]`)
	var match = re.FindString(text)
	if match != "" {
		divider = string(match[0])
		detectedNaming |= casing.Separated
		words = txt.Split(text, divider)
	}

	// Remove divider chars to analyze the core string
	var inputNoDivider = text
	if divider != "" {
		inputNoDivider = txt.Remove(text, divider)
	}

	if isAllLower(inputNoDivider) {
		detectedNaming = flag.TurnOff(detectedNaming, casing.Random)
		detectedNaming |= casing.Lower
		return detectedNaming, divider
	}
	if isAllUpper(inputNoDivider) {
		detectedNaming = flag.TurnOff(detectedNaming, casing.Random)
		detectedNaming |= casing.Upper
		return detectedNaming, divider
	}

	if len(words) == 1 {
		var runes = []rune(text)
		if unicode.IsLower(runes[0]) && containsUpper(runes[1:]) {
			detectedNaming = flag.TurnOff(detectedNaming, casing.Random)
			detectedNaming |= casing.Camel
		}
		if unicode.IsUpper(runes[0]) && containsUpper(runes[1:]) {
			detectedNaming = flag.TurnOff(detectedNaming, casing.Random)
			detectedNaming |= casing.Pascal
		}
		return detectedNaming, divider
	}

	// Check word-wise naming patterns
	var soFarCamel = isAllLower(words[0])
	var soFarPascal = isCapitalized(words[0])
	var soFarSentence = soFarPascal
	var soFarPing = isPing(words[0])
	var soFarPong = isPong(words[0])

	for _, word := range words[1:] {
		if !isAllLower(word) {
			soFarSentence = false
		}

		if !isCapitalized(word) {
			soFarCamel = false
			soFarPascal = false
		}

		if !isPing(word) {
			soFarPing = false
		}

		if !isPong(word) {
			soFarPong = false
		}

		if !soFarCamel && !soFarPascal && !soFarSentence && !soFarPing && !soFarPong {
			break
		}
	}

	if soFarCamel {
		detectedNaming = flag.TurnOff(detectedNaming, casing.Random)
		detectedNaming |= casing.Camel
	}
	if soFarPascal {
		detectedNaming = flag.TurnOff(detectedNaming, casing.Random)
		detectedNaming |= casing.Pascal
	}
	if soFarSentence {
		detectedNaming = flag.TurnOff(detectedNaming, casing.Random)
		detectedNaming |= casing.Sentence
	}
	if soFarPing {
		detectedNaming = flag.TurnOff(detectedNaming, casing.Random)
		detectedNaming |= casing.Pingpong
	}
	if soFarPong {
		detectedNaming = flag.TurnOff(detectedNaming, casing.Random)
		detectedNaming |= casing.Pongping
	}

	return detectedNaming, divider
}

// =================================================================
// private

func isAllLower(s string) bool {
	for _, r := range s {
		if !unicode.IsLower(r) {
			return false
		}
	}
	return true
}

func isAllUpper(s string) bool {
	for _, r := range s {
		if !unicode.IsUpper(r) {
			return false
		}
	}
	return true
}

func containsUpper(runes []rune) bool {
	for _, r := range runes { // no need for slices dependency
		if unicode.IsUpper(r) {
			return true
		}
	}
	return false
}

func isCapitalized(word string) bool {
	if word == "" {
		return false
	}
	var runes = []rune(word)
	return unicode.IsUpper(runes[0]) && isAllLower(string(runes[1:]))
}

func isPing(s string) bool {
	var isUpper = true
	for _, r := range s {
		if isUpper && !unicode.IsUpper(r) {
			return false
		}
		if !isUpper && !unicode.IsLower(r) {
			return false
		}
		isUpper = !isUpper
	}
	return true
}

func isPong(s string) bool {
	var isLower = true
	for _, r := range s {
		if isLower && !unicode.IsLower(r) {
			return false
		}
		if !isLower && !unicode.IsUpper(r) {
			return false
		}
		isLower = !isLower
	}
	return true
}
func capitalize(word string) string {
	if word == "" {
		return ""
	}
	var runes = []rune(word)
	return string(unicode.ToUpper(runes[0])) + txt.LowerCase(string(runes[1:]))
}
func addDivCamelPascal(text, div string) string {
	var result = txt.NewBuilder()
	var runes = []rune(text)
	for i := range runes {
		if i > 0 && unicode.IsUpper(runes[i]) && (i == len(runes)-1 || unicode.IsLower(runes[i+1])) {
			result.WriteText(div)
		}
		result.WriteSymbol(runes[i])
	}
	return result.ToText()
}
