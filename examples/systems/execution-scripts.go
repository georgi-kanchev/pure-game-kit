package example

import (
	"fmt"
	"pure-kit/engine/execution/script"
)

func Scripts() {
	var scr = script.New()

	scr.ExecuteFile("examples/data/script.lua")
	var result = scr.ExecuteFunction("Greet", "tosho")
	fmt.Printf("result: %v\n", result)
}
