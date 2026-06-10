// A very simple console-like command execution. Can act as an event system too.
//
// Turns a line of text (string command name + string parameters) into immediate mode custom code execution (no callbacks).
// Commands that are triggered at any time this frame can be received only throughout the entire next frame.
package command

import (
	"pure-game-kit/packages/internal"
	"pure-game-kit/packages/utility/number"
	"pure-game-kit/packages/utility/text"
)

// Command examples:
//
//	command_name: param0, param1, param2
//	log_messages: `hello, world!`, 123, true, false // strings are escaped by `
//	player_died: 35, 4 // from health points, by enemy index etc
//	toggle_fullscreen // parameters are optional
//	change_window_title: `My own window!`
//	debug: true
func Execute(command string) {
	command = text.Trim(text.Remove(command, "\r", "\n"))
	var replaced, originals = replaceStrings(command, quote, quote, placeholder)
	command = replaced

	var parts = text.Split(command, dividerParts)
	var name = text.Trim(substringUntilChar(parts[0], dividerName))
	parts[0] = text.Remove(parts[0], string(dividerName))
	parts[0] = text.Trim(text.Remove(parts[0], name))

	if len(parts) == 1 && parts[0] == "" {
		parts = []string{}
	}

	var originalStringIndex = 0
	for i := range parts {
		parts[i] = text.Trim(parts[i])

		if parts[i] == string(placeholder) {
			parts[i] = originals[originalStringIndex]
			originalStringIndex++
		}
	}
	internal.NewCommands[name] = parts
}

func JustExecuted(name string) bool {
	var _, has = internal.OldCommands[name]
	return has
}
func GrabText(name string, parameterIndex int) string {
	var params, _ = internal.OldCommands[name]
	if parameterIndex < 0 || parameterIndex >= len(params) {
		return ""
	}
	return params[parameterIndex]
}
func GrabNumber[T number.Number](name string, parameterIndex int) T {
	var params, _ = internal.OldCommands[name]
	if parameterIndex < 0 || parameterIndex >= len(params) {
		return T(0)
	}
	return text.ToNumber[T](params[parameterIndex])
}

// private ========================================================

const dividerParts = ","
const dividerName = ':'
const quote = '`'

func substringUntilChar(txt string, char rune) string {
	var index = text.IndexOf(txt, string(char))
	if index == -1 {
		return txt
	}
	return txt[:index]
}

const placeholder rune = '\x1A'

func replaceStrings(txt string, openQuote, closeQuote, placeholder rune) (replaced string, originals []string) {
	var replacedRunes []rune
	var extracted []string
	var currentString []rune
	var inString bool
	var runes = []rune(txt)
	for _, r := range runes {
		if !inString { // Found an opening quote
			if r == openQuote {
				inString = true
				replacedRunes = append(replacedRunes, placeholder)
			} else {
				replacedRunes = append(replacedRunes, r)
			}
		} else { // Found a closing quote
			if r == closeQuote {
				inString = false
				extracted = append(extracted, string(currentString))
				currentString = nil // Reset for the next string
			} else {
				currentString = append(currentString, r)
			}
		}
	}

	if inString { // if the user forgot the closing quote, we still save what we collected so it doesn't get lost
		extracted = append(extracted, string(currentString))
	}
	return string(replacedRunes), extracted
}
