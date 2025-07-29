package example

import (
	"fmt"
	"pure-kit/engine/execution/command"
)

func Commands() {
	command.Create("log_messages", func(parameters []string) (output string) {
		for i, v := range parameters {
			fmt.Printf("%v: %v\n", i, v)
		}
		return ""
	})

	command.Execute("log_messages: \"hello, world!\", test, 5")
}
