// A very simple terminal-like command creation and execution.
// Turns a line of text (string command name + string parameters) into custom code execution.
package command

import (
	"pure-game-kit/packages/utility/debug"
	"pure-game-kit/packages/utility/text"
)

func New(name string, execution func(parameters []string) (output string)) {
	commands[name] = execution
}

//=================================================================

// Command examples:
//
//	command_name: param0, param1, param2
//	log_messages: `hello, world!`, 1, 2, 3, 4, true, false, true
//	debug: true
//	change_window_title: `My own window!`
func Execute(command string) (output string) {
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

	var execution, has = commands[name]
	if !has {
		debug.LogError("Command not found: \"", command, "\"")
		return ""
	}

	var originalStringIndex = 0
	for i := range parts {
		parts[i] = text.Trim(parts[i])

		if parts[i] == string(placeholder) {
			parts[i] = originals[originalStringIndex]
			originalStringIndex++
		}
	}

	return execution(parts)
}

// private ========================================================

const dividerParts = ","
const dividerName = ':'
const quote = '"'

var commands = make(map[string]func([]string) string)

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
