package main

import (
	"pure-game-kit/packages/engine"
)

func main() {
	// var view = graphics.NewView(1)
	// var font = assets.LoadFont2("tools/sdf-font-generator/results/Montserrat-Medium.png",
	// "tools/sdf-font-generator/results/Montserrat-Medium.xml")
	// var obj = graphics.NewObject(0, 0)
	// obj.TextFont = font
	// obj.Text = "Hello, World!"

	// assets.LoadDefaultFont()

	engine.Run(60, func() {
		// var cycles = 50_000
		// res := 0
		// for i := 2; i < cycles; i++ {
		// 	isPrime := true
		// 	for j := 2; j*j <= i; j++ {
		// 		if i%j == 0 {
		// 			isPrime = false
		// 			break
		// 		}
		// 	}
		// 	if isPrime {
		// 		res++
		// 	}
		// }
	})
}
