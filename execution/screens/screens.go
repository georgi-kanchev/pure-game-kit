/*
A dead-simple screen manager. Essentially a global state machine with some
lifetime interface callback functions that update the screens.
Does not require manual updating since it is hooked to the engine pump.
Because of that, the main game loop may remain empty since each screen would be handled elsewhere.
*/
package screens

import (
	"pure-game-kit/debug"
	"pure-game-kit/internal"
)

type Screen interface {
	OnLoad()
	OnEnter()
	OnUpdate()
	OnExit()
}

//=================================================================

func Add(screen Screen, load bool) (screenId int) {
	if internal.CurrentScreen < 0 && internal.CurrentScreen >= len(internal.Screens) {
		debug.LogError("No screen found with id: ", screenId)
		return
	}

	internal.Screens = append(internal.Screens, screen)
	if load {
		screen.OnLoad()
	}
	return len(internal.Screens) - 1
}
func Enter(screenId int, load bool) {
	Current().OnExit()
	internal.CurrentScreen = screenId

	if load {
		Current().OnLoad()
	}
	Current().OnEnter()
}
func Reload() {
	for _, scr := range internal.Screens {
		scr.OnLoad()
	}
}

//=================================================================

func Current() Screen {
	return internal.Screens[internal.CurrentScreen]
}
func CurrentId() int {
	return internal.CurrentScreen
}
