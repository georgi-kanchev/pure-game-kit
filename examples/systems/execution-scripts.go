package example

import (
	"fmt"
	"pure-kit/engine/data/file"
	"pure-kit/engine/execution/script"
)

func Scripts() {
	var scr = script.New()
	var file = file.LoadText("examples/data/script.lua")
	scr.ExecuteCode(file)
	var result = scr.ExecuteFunction("Greet", "tosho")
	fmt.Printf("result: %v\n", result)
}
