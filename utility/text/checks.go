package text

import (
	"strings"
	"unicode"
	"unicode/utf8"
)

func Length(text string) int {
	return utf8.RuneCountInString(text)
}
func IndexOf(text, part string) int {
	return strings.Index(text, part)
}
func ContainsAll(text string, parts ...string) bool {
	for _, part := range parts {
		if !strings.Contains(text, part) {
			return false
		}
	}
	return true
}
func ContainsOneOf(text string, parts ...string) bool {
	for _, part := range parts {
		if strings.Contains(text, part) {
			return true
		}
	}
	return false
}
func CountOccurrences(text, part string) int {
	return strings.Count(text, part)
}

// Same as IsEmpty(...)
func IsBlank(text string) bool {
	return Trim(text) == ""
}

// Same as IsBlank(...)
func IsEmpty(text string) bool {
	return IsBlank(text)
}

func IsAllLetters(text string) bool {
	if len(text) == 0 {
		return false
	}
	for _, r := range text {
		if !unicode.IsLetter(r) {
			return false
		}
	}
	return true
}
func IsAllDigits(text string) bool {
	if len(text) == 0 {
		return false
	}
	for _, r := range text {
		if !unicode.IsDigit(r) {
			return false
		}
	}
	return true
}

func StartsWith(text, value string) bool {
	return strings.HasPrefix(text, value)
}
func EndsWith(text, value string) bool {
	return strings.HasSuffix(text, value)
}
