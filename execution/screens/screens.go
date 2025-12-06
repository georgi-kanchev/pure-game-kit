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
	screen.OnLoad()
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
