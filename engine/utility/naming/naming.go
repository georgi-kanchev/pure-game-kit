package naming

import (
	"math/rand"
	"regexp"
	"slices"
	"strings"
	"unicode"
)

type Naming int

const (
	RaNDomCasE Naming = 1 << iota
	lower
	UPPER
	camelCase
	PascalCase
	Sentence_case
	PiNgPoNg_CaSe
	pOnGpInG_cAsE
	Separated
)

func Apply(text string, naming Naming, divider string) string {
	if divider == "" {
		divider = ""
	}

	if naming == RaNDomCasE {
		var result strings.Builder
		var randBool = rand.Float64() < 0.5
		for _, r := range text {
			if randBool {
				result.WriteRune(unicode.ToLower(r))
			} else {
				result.WriteRune(unicode.ToUpper(r))
			}
		}
		return result.String()
	}

	detectedNaming, detectedDivider := Detect(text)
	var words = []string{text}
	if detectedDivider != "" {
		words = strings.Split(text, detectedDivider)
	}

	if len(words) == 1 &&
		divider != "" &&
		(hasFlag(detectedNaming, camelCase) || hasFlag(detectedNaming, PascalCase)) {
		words = strings.Split(addDivCamelPascal(words[0], divider), divider)
	}

	for i := range words {
		word := words[i]

		if hasFlag(naming, lower) {
			word = strings.ToLower(word)
		}

		if hasFlag(naming, UPPER) {
			word = strings.ToUpper(word)
		}

		if hasFlag(naming, camelCase) {
			if i == 0 {
				word = strings.ToLower(word)
			} else {
				word = capitalize(word)
			}
		}

		if hasFlag(naming, PascalCase) {
			word = capitalize(word)
		}

		if hasFlag(naming, Sentence_case) {
			if i == 0 {
				word = capitalize(word)
			} else {
				word = strings.ToLower(word)
			}
		}

		if hasFlag(naming, PiNgPoNg_CaSe) {
			var sb strings.Builder
			var isUpper = true
			for _, c := range word {
				if isUpper {
					sb.WriteRune(unicode.ToUpper(c))
				} else {
					sb.WriteRune(unicode.ToLower(c))
				}
				isUpper = !isUpper
			}
			word = sb.String()
		}

		if hasFlag(naming, pOnGpInG_cAsE) {
			var sb strings.Builder
			var isLower = true
			for _, c := range word {
				if isLower {
					sb.WriteRune(unicode.ToLower(c))
				} else {
					sb.WriteRune(unicode.ToUpper(c))
				}
				isLower = !isLower
			}
			word = sb.String()
		}

		words[i] = word
	}

	return strings.Join(words, divider)
}
func Detect(text string) (naming Naming, separator string) {
	if strings.TrimSpace(text) == "" {
		return RaNDomCasE, ""
	}

	var detectedNaming = RaNDomCasE
	var divider = ""
	var words = []string{text}
	var re = regexp.MustCompile(`[^a-zA-Z0-9]`)
	var match = re.FindString(text)
	if match != "" {
		divider = string(match[0])
		detectedNaming |= Separated
		words = strings.Split(text, divider)
	}

	// Remove divider chars to analyze the core string
	var inputNoDivider = text
	if divider != "" {
		inputNoDivider = strings.ReplaceAll(text, divider, "")
	}

	if isAllLower(inputNoDivider) {
		detectedNaming = removeFlag(detectedNaming, RaNDomCasE)
		detectedNaming |= lower
		return detectedNaming, divider
	}
	if isAllUpper(inputNoDivider) {
		detectedNaming = removeFlag(detectedNaming, RaNDomCasE)
		detectedNaming |= UPPER
		return detectedNaming, divider
	}

	if len(words) == 1 {
		var runes = []rune(text)
		if unicode.IsLower(runes[0]) && containsUpper(runes[1:]) {
			detectedNaming = removeFlag(detectedNaming, RaNDomCasE)
			detectedNaming |= camelCase
		}
		if unicode.IsUpper(runes[0]) && containsUpper(runes[1:]) {
			detectedNaming = removeFlag(detectedNaming, RaNDomCasE)
			detectedNaming |= PascalCase
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
		detectedNaming = removeFlag(detectedNaming, RaNDomCasE)
		detectedNaming |= camelCase
	}
	if soFarPascal {
		detectedNaming = removeFlag(detectedNaming, RaNDomCasE)
		detectedNaming |= PascalCase
	}
	if soFarSentence {
		detectedNaming = removeFlag(detectedNaming, RaNDomCasE)
		detectedNaming |= Sentence_case
	}
	if soFarPing {
		detectedNaming = removeFlag(detectedNaming, RaNDomCasE)
		detectedNaming |= PiNgPoNg_CaSe
	}
	if soFarPong {
		detectedNaming = removeFlag(detectedNaming, RaNDomCasE)
		detectedNaming |= pOnGpInG_cAsE
	}

	return detectedNaming, divider
}

// region private
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
	return slices.ContainsFunc(runes, unicode.IsUpper)
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
	return string(unicode.ToUpper(runes[0])) + strings.ToLower(string(runes[1:]))
}
func addDivCamelPascal(text, div string) string {
	var result strings.Builder
	var runes = []rune(text)
	for i := range runes {
		if i > 0 && unicode.IsUpper(runes[i]) && (i == len(runes)-1 || unicode.IsLower(runes[i+1])) {
			result.WriteString(div)
		}
		result.WriteRune(runes[i])
	}
	return result.String()
}
func hasFlag(value Naming, flag Naming) bool {
	return value&flag != 0
}
func removeFlag(value Naming, flag Naming) Naming {
	return value &^ flag
}

// endregion
