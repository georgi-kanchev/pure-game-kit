package internal

import "maps"

// this double buffer structure preserves commands for the next frame after the one they are executed in
var OldCommands, NewCommands map[string][]string = make(map[string][]string), make(map[string][]string)

func UpdateCommands() {
	clear(OldCommands)
	maps.Copy(OldCommands, NewCommands)
	clear(NewCommands)
}
