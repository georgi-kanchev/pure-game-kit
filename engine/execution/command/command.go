package command

import (
	"pure-kit/engine/internal"
	"strings"
)

// commands example
//		name: param0, param1, param2
// 		log_messages: "hello, world!", 1, 2, 3, 4, true, false, true
// 		debug: true
// 		change_window_title: "My own window!"

func Create(name string, execution func(parameters []string) (output string)) {
	commands[name] = execution
}

func Execute(command string) (output string) {
	command = strings.ReplaceAll(command, "\n", "")
	command = strings.ReplaceAll(command, "\r", "")
	command = strings.Trim(command, " ")
	var replaced, originals = internal.ReplaceQuotedStrings(command, quote, internal.Placeholder)
	command = replaced

	var parts = strings.Split(command, dividerParts)
	var name = strings.Trim(substringUntilChar(parts[0], dividerName), " ")
	parts[0] = strings.ReplaceAll(parts[0], string(dividerName), "")
	parts[0] = strings.Trim(strings.ReplaceAll(parts[0], name, ""), " ")

	if len(parts) == 1 && parts[0] == "" {
		parts = []string{}
	}

	var execution, has = commands[name]
	if !has {
		return ""
	}

	var originalStringIndex = 0
	for i := range parts {
		parts[i] = strings.Trim(parts[i], " ")

		if parts[i] == string(internal.Placeholder) {
			parts[i] = originals[originalStringIndex]
			originalStringIndex++
		}
	}

	return execution(parts)
}

// #region private

const dividerParts = ","
const dividerName = ':'
const quote = '"'

var commands = make(map[string]func([]string) string)

func substringUntilChar(text string, char rune) string {
	index := strings.IndexRune(text, char)
	if index == -1 {
		return text
	}
	return text[:index]
}

// #endregion
