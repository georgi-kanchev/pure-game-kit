package text

import (
	"strings"
	"unicode"
)

// Converts the text to lowercase
func ToLowerCase(text string) string {
	return strings.ToLower(text)
}

// Converts the text to UPPERCASE
func ToUpperCase(text string) string {
	return strings.ToUpper(text)
}

// Converts the text to camelCase
func ToCamelCase(text string) string {
	words := SplitWords(text)
	for i := 1; i < len(words); i++ {
		words[i] = ToTitleCase(words[i])
	}
	return strings.Join(words, "")
}

// Converts the text to PascalCase
func ToPascalCase(text string) string {
	words := SplitWords(text)
	for i := range words {
		words[i] = ToTitleCase(words[i])
	}
	return strings.Join(words, "")
}

// Converts the text to Sentence case
func ToSentenceCase(text string) string {
	if len(text) == 0 {
		return ""
	}
	res := strings.ToLower(text)
	return strings.ToUpper(string(res[0])) + res[1:]
}

// Converts the text to PiNgPoNg CaSe (UPPERCASE starts first)
func ToPingPongCase(text string) string {
	var sb strings.Builder
	upper := true
	for _, r := range text {
		if unicode.IsLetter(r) {
			if upper {
				sb.WriteRune(unicode.ToUpper(r))
			} else {
				sb.WriteRune(unicode.ToLower(r))
			}
			upper = !upper
		} else {
			sb.WriteRune(r)
		}
	}
	return sb.String()
}

// Converts the text to pOnGpInG cAsE (lowercase starts first)
func ToPongPingCase(text string) string {
	var sb strings.Builder
	upper := false
	for _, r := range text {
		if unicode.IsLetter(r) {
			if upper {
				sb.WriteRune(unicode.ToUpper(r))
			} else {
				sb.WriteRune(unicode.ToLower(r))
			}
			upper = !upper
		} else {
			sb.WriteRune(r)
		}
	}
	return sb.String()
}

// Converts the text to kebab-case
func ToKebabCase(text string) string {
	return strings.Join(SplitWords(text), "-")
}

// Converts the text to snake_case
func ToSnakeCase(text string) string {
	return strings.Join(SplitWords(text), "_")
}

// Converts the text to Title Case
func ToTitleCase(text string) string {
	var sb strings.Builder
	sb.Grow(len(text)) // pre-allocate memory for performance
	isAtWordStart := true

	for _, r := range text {
		if isSeparator(r) {
			sb.WriteRune(r)
			isAtWordStart = true
		} else {
			if isAtWordStart {
				sb.WriteRune(unicode.ToUpper(r))
				isAtWordStart = false
			} else {
				sb.WriteRune(unicode.ToLower(r))
			}
		}
	}

	return sb.String()
}
