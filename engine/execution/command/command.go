package command

import (
	"pure-kit/engine/internal"
	"pure-kit/engine/utility/text"
)

func New(name string, execution func(parameters []string) (output string)) {
	commands[name] = execution
}

//=================================================================

// example
//
//	"command_name: param0, param1, param2"
//	"log_messages: `hello, world!`, 1, 2, 3, 4, true, false, true"
//	"debug: true"
//	"change_window_title: `My own window!`"
func Execute(command string) (output string) {
	command = text.Trim(text.Remove(command, "\r", "\n"))
	var replaced, originals = internal.ReplaceQuotedStrings(command, quote, internal.Placeholder)
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
		return ""
	}

	var originalStringIndex = 0
	for i := range parts {
		parts[i] = text.Trim(parts[i])

		if parts[i] == string(internal.Placeholder) {
			parts[i] = originals[originalStringIndex]
			originalStringIndex++
		}
	}

	return execution(parts)
}

//=================================================================
// private

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
