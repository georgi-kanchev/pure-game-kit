package main

import (
	"pure-game-kit/debug"
	example "pure-game-kit/examples/systems"
)

func main() {
	debug.ProfileCPU(5)
	example.Shapes()
}
