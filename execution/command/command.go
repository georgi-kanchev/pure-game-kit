/*
A very simple terminal-like command creation and execution.
Turns a line of text (string command name + string parameters) into custom code execution.
*/
package command

import (
	"pure-game-kit/debug"
	"pure-game-kit/internal"
	"pure-game-kit/utility/text"
)

func New(name string, execution func(parameters []string) (output string)) {
	commands[name] = execution
}

//=================================================================

/*
Command examples:

	command_name: param0, param1, param2
	log_messages: `hello, world!`, 1, 2, 3, 4, true, false, true
	debug: true
	change_window_title: `My own window!`
*/
func Execute(command string) (output string) {
	command = text.Trim(text.Remove(command, "\r", "\n"))
	var replaced, originals = internal.ReplaceStrings(command, quote, quote, internal.Placeholder)
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
