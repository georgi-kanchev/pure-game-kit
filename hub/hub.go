package main

import (
	"pure-kit/engine/utility/time"
	"pure-kit/engine/window"
)

func main() {
	for window.KeepOpen() {
		time.Update()
	}
}
