package main

import (
	"fmt"
	"pure-game-kit/packages/engine"
	"pure-game-kit/packages/input/keyboard"
)

func main() {
	engine.Initialize("pure-game-kit", 60, 0, false, false)

	// var view = graphics.NewView(1)
	// var font = assets.LoadFont2("tools/sdf-font-generator/results/Montserrat-Medium.png",
	// "tools/sdf-font-generator/results/Montserrat-Medium.xml")
	// var obj = graphics.NewObject(0, 0)
	// obj.TextFont = font
	// obj.Text = "Hello, World!"

	// assets.LoadDefaultFont()

	engine.Run(func() {
		fmt.Printf("keyboard.Input(): %v\n", keyboard.Input())
	})
}
