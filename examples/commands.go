package example

import (
	"fmt"
	"pure-kit/engine/execution/flow"
)

func Commands() {
	flow.CommandCreate("log_messages", func(parameters []string) (output string) {
		for i, v := range parameters {
			fmt.Printf("%v: %v\n", i, v)
		}
		return ""
	})

	flow.CommandExecute("log_messages: \"hello, world!\", test, 5")
}
