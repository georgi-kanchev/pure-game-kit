/*
A dead-simple screen/scene manager. Essentially a global state machine with some
scene-specific interface callback functions that always updates the current screen/scene.
Does not require manual updating since it is hooked to the engine pump.
Because of that, the main game loop may remain empty since each screen would be handled elsewhere.
*/
package screens

import "pure-game-kit/internal"

type Screen interface {
	OnLoad()
	OnEnter()
	OnUpdate()
	OnExit()
}

func Add(screen Screen, load bool) (screenId int) {
	internal.Screens = append(internal.Screens, screen)
	if load {
		screen.OnLoad()
	}
	return len(internal.Screens) - 1
}
func Enter(screenId int, load bool) {
	internal.Screens[internal.CurrentScreen].OnExit()
	internal.CurrentScreen = screenId

	if load {
		internal.Screens[internal.CurrentScreen].OnLoad()
	}
	internal.Screens[internal.CurrentScreen].OnEnter()
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
