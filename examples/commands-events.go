package example

import (
	"pure-game-kit/packages/execution/command"
	"pure-game-kit/packages/input/keyboard"
	"pure-game-kit/packages/input/keyboard/key"
	"pure-game-kit/packages/utility/debug"
	"pure-game-kit/packages/utility/text"
	"pure-game-kit/packages/utility/time"
	"pure-game-kit/packages/window"
)

func CommandsEvents() {
	window.Create("example - commands/events", true, true)

	var counter = 0
	for window.KeepOpen() {
		if keyboard.IsKeyJustPressed(key.A) {
			counter++
			command.Execute(text.New("print_runtime: ", counter))
		}

		outsideScope()
	}
}

// private ========================================================

func outsideScope() {
	if command.JustExecuted("print_runtime") {
		var counter = command.GrabNumber[int]("print_runtime", 0)
		debug.Print(time.Running(), " counter: ", counter)
	}
}
